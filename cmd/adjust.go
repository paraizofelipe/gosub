package cmd

import (
	"flag"

	"github.com/paraizofelipe/gosub/srt"
)

type AdjustCmd struct {
	file  string
	ms    int
	flags *flag.FlagSet
}

func NewAdjustCmd() *AdjustCmd {
	c := &AdjustCmd{
		flags: flag.NewFlagSet("adjust", flag.ContinueOnError),
	}
	c.flags.StringVar(&c.file, "file", "", "path of subtitle file")
	c.flags.IntVar(&c.ms, "ms", 0, "aditional milliseconds")

	return c
}

func (_c *AdjustCmd) Name() string {
	return _c.flags.Name()
}

func (_c *AdjustCmd) Init(args []string) error {
	return _c.flags.Parse(args)
}

func (_c *AdjustCmd) Run() (err error) {
	subSrt := srt.NewSubSrt()

	_, err = subSrt.Read(_c.file)
	if err != nil {
		return
	}

	err = subSrt.AdjustTime(_c.ms)
	if err != nil {
		return
	}

	err = subSrt.Write(_c.file)
	if err != nil {
		return
	}
	return
}
