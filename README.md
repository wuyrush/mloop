# `mloop` - Loop playing music

## What?

I am fed up with my noisy neighbors. So I wrote this small tool to tease them.

It loops playing audio files(for now mp3 is the only supported codec) in the given directory for given amount of time, starting at a given time in the current date.

## Usage
```shell
mloop -dir <audio-file-dir> -d <duration> [-s <play-start-time>] [-v]
```

Example:
```shell
# play all audio files in current directory now for 30min
mloop -d 30m
# play audio files in directory ~/Download/music for 1 hour, starting from 19:00 today
mloop -dir ~/Download/music -d 1h -s 19:00
# ctrl-c to exit mloop if needed 
```

## Build

First make sure Go is installed on your box. If not, [get it](https://golang.org/doc/install).

Then build and enjoy:
```shell
git clone https://github.com/wuyrush/mloop

cd mloop
go build github.com/wuyrush/mloop

./mloop --help
Usage of mloop:
  -d duration
    	Amount of time to loop playing the given audio files. Default to 0
  -dir string
    	Path to audio file directory. Default to current directory (default ".")
  -s value
    	start time in HH:MM format. Default to current time (default 19:00)
  -v	Turn on verbose mode
```

## TODOs 

### 08/17/19
* Accept single audio file as input as well
* Support more codecs
* Consider playable audio files only; For now `mloop` will complain and exit if it processes a file it cannot play
* The only play order supported by `mloop` is ascending lexico order based on audio file filenames; maybe add an option for random order?

