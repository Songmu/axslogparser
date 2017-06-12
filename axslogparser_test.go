package axslogparser

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

var loc = func() *time.Location {
	t, _ := time.Parse(clfTimeLayout, "11/Jun/2017:05:56:04 +0900")
	return t.Location()
}()

func pfloat64(f float64) *float64 {
	return &f
}

var parseTests = []struct {
	Name   string
	Input  string
	Output Log
}{
	{
		Name:  "combined",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200 741 "-" "mackerel-http-checker/0.0.1" "-"`,
		Output: Log{
			Host:    "10.0.0.11",
			User:    "-",
			Time:    time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request: "GET / HTTP/1.1",
			Status:  200,
			Size:    741,
			Referer: "-",
			UA:      "mackerel-http-checker/0.0.1",
		},
	},
	{
		Name:  "common",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200 741`,
		Output: Log{
			Host:    "10.0.0.11",
			User:    "-",
			Time:    time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request: "GET / HTTP/1.1",
			Status:  200,
			Size:    741,
		},
	},
	{
		Name:  "common with empty response",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 204 -`,
		Output: Log{
			Host:    "10.0.0.11",
			User:    "-",
			Time:    time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request: "GET / HTTP/1.1",
			Status:  204,
			Size:    0,
		},
	},
	{
		Name:  "common with vhost",
		Input: `log.example.com 10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 404 741`,
		Output: Log{
			VirtualHost: "log.example.com",
			Host:        "10.0.0.11",
			User:        "-",
			Time:        time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:     "GET / HTTP/1.1",
			Status:      404,
			Size:        741,
		},
	},
	{
		Name:  "unescape",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET /?foo=bar HTTP/1.1" 200 741 "\\\thoge" "UA \"fake\""`,
		Output: Log{
			Host:    "10.0.0.11",
			User:    "-",
			Time:    time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request: "GET /?foo=bar HTTP/1.1",
			Status:  200,
			Size:    741,
			Referer: "\\\thoge",
			UA:      `UA "fake"`,
		},
	},
	{
		Name: "ltsv",
		Input: "time:08/Mar/2017:14:12:40 +0900\t" +
			"host:192.0.2.1\t" +
			"req:POST /api/v0/tsdb HTTP/1.1\t" +
			"status:200\t" +
			"size:36\t" +
			"ua:mackerel-agent/0.31.2 (Revision 775fad2)\t" +
			"reqtime:0.087\t" +
			"taken_sec:0.087\t" +
			"vhost:mackerel.io",
		Output: Log{
			VirtualHost: "mackerel.io",
			Host:        "192.0.2.1",
			Time:        time.Date(2017, time.March, 8, 14, 12, 40, 0, loc),
			TimeStr:     "08/Mar/2017:14:12:40 +0900",
			Request:     "POST /api/v0/tsdb HTTP/1.1",
			Status:      200,
			Size:        36,
			UA:          "mackerel-agent/0.31.2 (Revision 775fad2)",
			ReqTime:     pfloat64(0.087),
			TakenSec:    pfloat64(0.087),
		},
	},
	{
		Name:  "unescape(trailing space after escaped double quote) (TODO)",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET /?foo=bar HTTP/1.1" 200 741 "\" "`,
		Output: Log{
			Host:    "10.0.0.11",
			User:    "-",
			Time:    time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request: "GET /?foo=bar HTTP/1.1",
			Status:  200,
			Size:    741,
			Referer: `" `,
		},
	},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		t.Logf("testing: %s\n", tt.Name)
		if strings.Contains(tt.Name, "(TODO)") {
			t.Skipf("skip test: %s", tt.Name)
		}
		l, err := Parse(tt.Input)
		if err != nil {
			t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
			continue
		}
		if !reflect.DeepEqual(*l, tt.Output) {
			t.Errorf("%s(parse): \n out =%+v\n want %+v", tt.Name, *l, tt.Output)
		}
	}
}
