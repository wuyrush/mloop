package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type timepoint time.Time

func (tp *timepoint) String() string {
	return time.Time(*tp).String()
}

func (tp *timepoint) Set(s string) error {
	t, err := time.Parse("15:04", s)
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
		ad    = flag.String("dir", "", "Path to audio file directory")
		d     = flag.Duration("duration", 0*time.Second, "Loop playing the given audio files for how long?")
		start timepoint
	)
	flag.Var(&start, "start", "start time in HH:MM format")

	flag.Parse()
	log.Infof("Loop playing audio files in %s for %s, starting at %s", *ad, *d, &start)

	timer := time.NewTimer(*d)
	// spawn a goroutine to 1) loop play music and 2) listen to timer signal for exit
	// when the goroutine exit, inform the waiting main goroutine to exit as well
	loop(*ad, timer.C)
}

func loop(dir string, _ <-chan time.Time) {
	// assume dir points to an audio file, testing beep
	clog := log.WithField("file", dir)
	f, err := os.Open(dir)
	if err != nil {
		clog.WithError(err).Error("error opening file")
		return
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		clog.WithError(err).Errorf("error decoding file")
		return
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
	select {}
}

func iteratorFunc(dir string) func() beep.Streamer {
	// get paths of all direct descendants of dir
	paths := make([]string)
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		// skip if p points to a directory
		if info.IsDir() {
			return filepath.SkipDir
		}
		paths = append(paths, p)
		return nil
	})
	pointer := 0
	return func() beep.Streamer {
		// generate a streamer via Seq(), which compose the real streamer and
		// a dummy streamer(beep.Callback) to close the real streamer together
	}
}

/*
1. walk all the files of the most shallow level in the given directory
2. Have a pointer to point to the file being played currently
3. have a function to loop generating a streamer based on the file pointer
4. pass the function to beep.Iterate() so that it produces a looping streamer
5. play the streamer we got from 4).
*/
