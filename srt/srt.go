package srt

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/paraizofelipe/gosub/models"
)

type SubSrt struct {
	Subtitles   []models.Srt
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

// Read --
func (_s *SubSrt) Read(filePath string) (srtList []models.Srt, err error) {
	var (
		index int
		srt   models.Srt
	)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {

		line := scanner.Text()
		if _s.isTimeLine(line) {
			srt.TimeLine = line
			continue
		}
		if _s.isIndexLine(line) {
			if index, err = strconv.Atoi(line); err != nil {
				return
			}
			srt.IndexLine = index
			continue
		}
		if line != "" {
			srt.TextLines = append(srt.TextLines, line)
			continue
		}

		_s.Subtitles = append(_s.Subtitles, srt)
		srt = models.Srt{}
	}
	return
}

// ToString --
func (_c SubSrt) ToString() (strSrt string, err error) {
	for _, srt := range _c.Subtitles {
		strSrt += fmt.Sprint(srt)
	}
	return
}

// Write --
func (_s SubSrt) Write(filePath string) (err error) {
	var strSrt string

	if strSrt, err = _s.ToString(); err != nil {
		return
	}

	out := []byte(strSrt)
	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// AdjustTime --
func (_s *SubSrt) AdjustTime(addTime int, indexRange [2]int) (err error) {

	for index, srt := range _s.Subtitles {
		if (index >= indexRange[0]-1 && index <= indexRange[1]-1) || (indexRange[0] == 0 && indexRange[1] == 0) {
			start := fmt.Sprint(_s.reTimeLine.ReplaceAllString(srt.TimeLine, "${start}"))
			end := fmt.Sprint(_s.reTimeLine.ReplaceAllString(srt.TimeLine, "${end}"))

			startTime, err := _s.addTime(start, addTime)
			if err != nil {
				return err
			}
			endTime, err := _s.addTime(end, addTime)
			if err != nil {
				return err
			}

			start = _s.timeToSubtitleTime(startTime)
			end = _s.timeToSubtitleTime(endTime)

			_s.Subtitles[index].TimeLine = fmt.Sprintf("%s --> %s", start, end)
		}
	}
	return
}

// Get --
func (_s SubSrt) Get(index int) (srt models.Srt) {
	return _s.Subtitles[index-1]
}
