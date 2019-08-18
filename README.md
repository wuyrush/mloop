# `mloop` - Loop playing music

## What?

I am fed up with my noisy neighbors. So I wrote this small tool to tease them.

It loops playing audio files(for now mp3 is the only supported codec) in the given directory for given amount of time, starting at a given time in the current date.

## Usage

```shell
mloop -dir <audio-file-dir> -d 15m -s <play-start-time> [-v]
```

## Build

First make sure Go is installed on your box. If not, [get it](https://golang.org/doc/install).

Then build and enjoy:
```shell
git clone https://github.com/wuyrush/mloop

cd mloop
go build github.com/wuyrush/mloop

./mloop --help
Usage of ./mloop:
  -d duration
    	Amount of time to loop playing the given audio files
  -dir string
    	Path to audio file directory
  -s value
    	start time in HH:MM format (default 18:10)
  -v	Turn on verbose mode
```

