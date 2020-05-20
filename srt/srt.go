package srt

import (
	"bufio"
	"fmt"
	"gosub/models"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SubSrt struct {
	reTimeLine  *regexp.Regexp
	reIndexLine *regexp.Regexp
}

func NewSubSrt() *SubSrt {
	timePattern := `(?P<start>(\d{2}:\d{2}:\d+\,\d+))( \-\-\> )(?P<end>(\d{2}:\d{2}:\d+\,\d+))`
	indexPattern := `^\d+$`

	return &SubSrt{
		reTimeLine:  regexp.MustCompile(timePattern),
		reIndexLine: regexp.MustCompile(indexPattern),
	}
}

// isIndexLine --
func (_s SubSrt) isIndexLine(line string) bool {
	return _s.reIndexLine.MatchString(line)
}

// isTimeLine --
func (_s SubSrt) isTimeLine(line string) bool {
	return _s.reTimeLine.MatchString(line)
}

// formatTime --
func (_s SubSrt) formatTime(time string) string {
	s := strings.Split(time, ",")
	return fmt.Sprintf("2006-01-02T%s.%sZ", s[0], s[1])
}

// timeToSubtitleTime --
func (_s SubSrt) timeToSubtitleTime(t time.Time) string {
	return fmt.Sprintf("%.2d:%.2d:%.2d,%.3d", t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)
}

// AddTime ---
func (_s SubSrt) addTime(srcTime string, addTime int) (timeOffset time.Time, err error) {
	t, err := time.Parse(time.RFC3339, _s.formatTime(srcTime))
	if err != nil {
		return
	}
	timeOffset = t.Add(time.Millisecond * time.Duration(addTime))
	return
}

// Reader --
func (_s SubSrt) Reader(filePath string) (srt models.Srt, err error) {
	var (
		index     int
		textLines []string
	)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if _s.isTimeLine(line) {
			srt.TimeLines = append(srt.TimeLines, line)
			continue
		}
		if _s.isIndexLine(line) {
			if index, err = strconv.Atoi(line); err != nil {
				return
			}
			srt.IndexLines = append(srt.IndexLines, index)
			continue
		}
		if line != "" {
			textLines = append(textLines, line)
			continue
		}

		srt.TextLines = append(srt.TextLines, textLines)
		textLines = []string{}
	}
	return
}

// Writer --
func (_s SubSrt) Writer(filePath string, srt models.Srt) (err error) {
	var out []byte

	for index := 0; index < len(srt.IndexLines); index++ {
		i := strconv.Itoa(srt.IndexLines[index])
		out = append(out, []byte(i+"\n")...)
		out = append(out, []byte(srt.TimeLines[index]+"\n")...)
		for _, text := range srt.TextLines[index] {
			out = append(out, []byte(text+"\n")...)
		}
		out = append(out, []byte("\n")...)
	}

	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// AdjustTime --
func (_s SubSrt) AdjustTime(srt models.Srt, addTime int) (models.Srt, error) {

	for index := range srt.TimeLines {
		start := fmt.Sprint(_s.reTimeLine.ReplaceAllString(srt.TimeLines[index], "${start}"))
		end := fmt.Sprint(_s.reTimeLine.ReplaceAllString(srt.TimeLines[index], "${end}"))

		startTime, err := _s.addTime(start, addTime)
		if err != nil {
			return srt, err
		}
		endTime, err := _s.addTime(end, addTime)
		if err != nil {
			return srt, err
		}

		start = _s.timeToSubtitleTime(startTime)
		end = _s.timeToSubtitleTime(endTime)

		srt.TimeLines[index] = fmt.Sprintf("%s --> %s", start, end)
	}
	return srt, nil
}
