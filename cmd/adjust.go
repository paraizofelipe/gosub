package cmd

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/paraizofelipe/gosub/srt"
)

type IndexRange [2]int

func (_i *IndexRange) Set(value string) error {
	var err error
	s := strings.Split(value, "-")
	for index := range s {
		v, err := strconv.Atoi(s[index])
		if err != nil {
			return err
		}
		_i[index] = v
	}
	return err
}

func (_i IndexRange) String() string {
	var str string
	for _, v := range _i {
		str += strconv.Itoa(v)
	}
	return str
}

type AdjustCmd struct {
	file       string
	ms         int
	indexRange IndexRange
	flags      *flag.FlagSet
}

func NewAdjustCmd() *AdjustCmd {
	c := &AdjustCmd{
		ms:    -1,
		flags: flag.NewFlagSet("adjust", flag.ContinueOnError),
	}

	c.flags.StringVar(&c.file, "file", "", "path of subtitle file")
	c.flags.IntVar(&c.ms, "ms", 0, "aditional milliseconds")
	c.flags.Var(&c.indexRange, "range", "initial and final index values separated by a \"-\" character exp: 8-21")

	return c
}

func (_c *AdjustCmd) Name() string {
	return _c.flags.Name()
}

func (_c *AdjustCmd) Init(args []string) error {
	return _c.flags.Parse(args)
}

func (_c AdjustCmd) ValidateFlags() (err error) {
	if _c.ms == -1 {
		return errors.New("value to ms is required!")
	}
	if _c.indexRange[0] > _c.indexRange[1] {
		return errors.New("invalid index values!")
	}
	return
}

func (_c *AdjustCmd) Run() (err error) {
	subSrt := srt.NewSubSrt()

	_, err = subSrt.Read(_c.file)
	if err != nil {
		return
	}

	err = subSrt.AdjustTime(_c.ms, _c.indexRange)
	if err != nil {
		return
	}

	err = subSrt.Write(_c.file)
	if err != nil {
		return
	}
	return
}
