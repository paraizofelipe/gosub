# gosub

A simple application to adjust the time in subtitle files .srt.  

## Prerequisites

- golang >= 1.14.3

## Installing

```bash=
$ go get -u github.com/paraizofelipe/gosub
```

## Run

```bash=
$ gosub adjust -file example.s01e02.srt -ms -1000
```

## Running the tests

```bash=
$ make test
```
