package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/paraizofelipe/gosub/srt"
)

func main() {
	subSrt := srt.NewSubSrt()

	flag.Parse()
	fileName := flag.Arg(0)
	strMs := flag.Arg(1)

	ms, err := strconv.Atoi(strMs)
	if err != nil {
		log.Fatal(err)
	}

	srt, err := subSrt.Reader(fileName)
	if err != nil {
		log.Fatal(err)
	}

	srt, err = subSrt.AdjustTime(srt, ms)
	if err != nil {
		log.Fatal(err)
	}

	err = subSrt.Writer(fileName, srt)
	if err != nil {
		log.Fatal(err)
	}
}
