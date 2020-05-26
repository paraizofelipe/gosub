package cmd

import (
	"flag"

	"github.com/paraizofelipe/gosub/srt"
)

type ShowCmd struct {
	flags *flag.FlagSet
	index int
	file  string
}

func NewShowCmd() *ShowCmd {
	c := &ShowCmd{
		flags: flag.NewFlagSet("get", flag.ContinueOnError),
	}
	c.flags.StringVar(&c.file, "file", "", "path of subtitle file")
	c.flags.IntVar(&c.index, "index", 0, "subtitle position index")

	return c
}

func (_c *ShowCmd) Name() string {
	return _c.flags.Name()
}

func (_c *ShowCmd) Init(args []string) error {
	return _c.flags.Parse(args)
}

func (_c *ShowCmd) Run() (err error) {
	subSrt := srt.NewSubSrt()

	_, err = subSrt.Read(_c.file)
	if err != nil {
		return
	}

	return
}
