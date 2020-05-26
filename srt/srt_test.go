package srt

import (
	"log"
	"testing"
	"time"

	"github.com/paraizofelipe/gosub/models"

	"github.com/stretchr/testify/assert"
)

func getTime(strTime string) (t time.Time) {
	t, err := time.Parse(time.RFC3339, strTime)
	if err != nil {
		log.Println(err)
	}
	return
}

func TestFormatTime(t *testing.T) {
	cases := []struct {
		description string
		in          string
		expect      string
	}{
		{
			description: "simple test with zero in front",
			in:          "01:01:01,001",
			expect:      "2006-01-02T01:01:01.001Z",
		},
		{
			description: "simple test",
			in:          "10:10:10,100",
			expect:      "2006-01-02T10:10:10.100Z",
		},
	}

	subSrt := NewSubSrt()

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			time := subSrt.formatTime(test.in)
			assert.Equal(t, time, test.expect)
		})
	}
}

func TestTimeToSubtitleTime(t *testing.T) {
	cases := []struct {
		description string
		in          time.Time
		expect      string
	}{
		{
			description: "simple test time",
			in:          getTime("2020-11-10T23:10:20.00Z"),
			expect:      "23:10:20,000",
		},
		{
			description: "simple test with milliseconds",
			in:          getTime("2006-01-02T23:10:20.123Z"),
			expect:      "23:10:20,123",
		},
	}

	subSrt := NewSubSrt()

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			time := subSrt.timeToSubtitleTime(test.in)
			assert.Equal(t, test.expect, time)
		})
	}
}

func TestAddTime(t *testing.T) {
	type params struct {
		time string
		ms   int
	}

	cases := []struct {
		description string
		in          params
		expect      time.Time
	}{
		{
			description: "add 1 second",
			in:          params{"23:10:20,000", 1000},
			expect:      getTime("2006-01-02T23:10:21.000Z"),
		},
		{
			description: "add 100 milliseconds",
			in:          params{"23:10:20,000", 100},
			expect:      getTime("2006-01-02T23:10:20.100Z"),
		},
	}

	subSrt := NewSubSrt()

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			time, err := subSrt.addTime(test.in.time, test.in.ms)
			assert.NoError(t, err)
			assert.Equal(t, test.expect, time)
		})
	}

}

func TestAdjustTime(t *testing.T) {
	cases := []struct {
		description string
		srt         []models.Srt
		in          int
		indexRange  [2]int
		expect      []models.Srt
	}{
		{
			description: "add 1 second",
			indexRange:  [2]int{0, 0},
			srt: []models.Srt{
				{TimeLine: "00:00:01,000 --> 00:00:02,000"},
			},
			in: 1000,
			expect: []models.Srt{
				{TimeLine: "00:00:02,000 --> 00:00:03,000"},
			},
		},
		{
			description: "add 1 milliseconds",
			indexRange:  [2]int{0, 0},
			srt: []models.Srt{
				{TimeLine: "00:00:02,100 --> 00:00:03,100"},
			},
			in: 100,
			expect: []models.Srt{
				{TimeLine: "00:00:02,200 --> 00:00:03,200"},
			},
		},
		{
			description: "add 1 second a range 1-4",
			indexRange:  [2]int{1, 3},
			srt: []models.Srt{
				{TimeLine: "00:00:05,100 --> 00:00:06,100"},
				{TimeLine: "00:00:07,100 --> 00:00:08,100"},
				{TimeLine: "00:00:09,100 --> 00:00:10,100"},
			},
			in: 1000,
			expect: []models.Srt{
				{TimeLine: "00:00:06,100 --> 00:00:07,100"},
				{TimeLine: "00:00:08,100 --> 00:00:09,100"},
				{TimeLine: "00:00:10,100 --> 00:00:11,100"},
			},
		},
	}

	for _, test := range cases {
		subSrt := NewSubSrt()
		subSrt.Subtitles = test.srt

		t.Run(test.description, func(t *testing.T) {
			err := subSrt.AdjustTime(test.in, test.indexRange)

			assert.NoError(t, err)
			for index, expected := range test.expect {
				assert.Equal(t, expected, subSrt.Subtitles[index])
			}
		})
	}
}
