package models

import (
	"fmt"
	"strings"
)

type Srt struct {
	IndexLine int
	TimeLine  string
	TextLines []string
}

func (_m Srt) String() string {
	return fmt.Sprintf("%d\n%s\n%s\n\n", _m.IndexLine, _m.TimeLine, strings.Join(_m.TextLines, "\n"))
}
