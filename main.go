package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	log "github.com/sirupsen/logrus"
)

const (
	ResampleQualifyIdx = 4
	SpeakerSampleRate  = beep.SampleRate(48000)
	StartTimeFormat    = "15:04"
)

func init() {
}

type timepoint time.Time

func (tp *timepoint) String() string {
	return time.Time(*tp).Format(StartTimeFormat)
}

func (tp *timepoint) Set(s string) error {
	t, err := time.Parse(StartTimeFormat, s)
	if err != nil {
		return err
	}
	h, m, sec := t.Clock()
	now := time.Now()
	re := time.Date(now.Year(), now.Month(), now.Day(), h, m, sec, 0, time.Local)
	if re.Before(now) {
		return fmt.Errorf("Start time %v is passed", re)
	}
	*tp = timepoint(re)
	return nil
}

func main() {
	var (
		ad      = flag.String("dir", "", "Path to audio file directory")
		d       = flag.Duration("d", 0*time.Second, "Amount of time to loop playing the given audio files")
		start   = timepoint(time.Now())
		verbose = flag.Bool("v", false, "Turn on verbose mode")
	)
	flag.Var(&start, "s", "start time in HH:MM format")
	flag.Parse()
	// setup logging
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	// get the absolute filepath of audio file dir
	dir, err := filepath.Abs(*ad)
	if err != nil {
		log.WithError(err).Fatal("failed to get the absolute path of audio file directory")
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("Loop playing audio files in %s for %s, starting at %s", dir, *d, &start)
	if time.Now().Before(time.Time(start)) {
		log.Infof("Wait till %s", &start)
		select {
		case sig := <-sigChan:
			log.WithField("signal", sig).Info("got exit signal while waiting. Exit")
			return
		case <-time.After(time.Time(start).Sub(time.Now())):
		}
	}
	// initialize speaker. Note beep.speaker has fixed sample rate
	speaker.Init(SpeakerSampleRate, SpeakerSampleRate.N(time.Second/10))
	timer := time.NewTimer(*d)
	exit := make(chan struct{})
	var wg sync.WaitGroup
	f := loopFunc(dir, exit)
	wg.Add(1)
	go f(&wg)

	select {
	case sig := <-sigChan:
		log.WithField("signal", sig).Info("got exit signal")
	case <-timer.C:
		log.Info("play time's up")
	}
	// notify the looping goroutine to exit and
	close(exit)
	wg.Wait()
	log.Info("main exits")
}

func loopFunc(dir string, exit <-chan struct{}) func(*sync.WaitGroup) {
	// get paths of all direct descendants of dir
	var paths []string
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		log.Debugf("walk %s", p)
		// skip if p points to a directory
		if p == dir {
			return nil
		} else if info.IsDir() {
			return filepath.SkipDir
		}
		paths = append(paths, p)
		return nil
	})
	log.WithField("files", paths).Debugf("done enumerating files in directory %s", dir)
	idx, done := 0, make(chan bool)
	return func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			p := paths[idx]
			cont := func() bool {
				clog := log.WithField("filepath", p)
				f, err := os.Open(p)
				if err != nil {
					clog.WithError(err).Fatal("error open file")
				}
				streamer, format, err := mp3.Decode(f)
				if err != nil {
					clog.WithError(err).Fatal("error decoding file with mp3 codec")
				}
				defer streamer.Close()

				resampled := beep.Resample(ResampleQualifyIdx, format.SampleRate, SpeakerSampleRate, streamer)
				speaker.Play(beep.Seq(resampled, beep.Callback(func() {
					select {
					case <-exit:
					case done <- true:
					}
				})))
				// wait till either the current song playback finishes or exit signal comes
				select {
				case <-exit:
					return false
				case <-done:
					return true
				}
			}()
			if !cont {
				return
			}
			idx = (idx + 1) % len(paths)
		}
	}
}
