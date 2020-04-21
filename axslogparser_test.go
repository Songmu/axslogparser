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
		Name:  "[Apache] combined",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200 741 "-" "mackerel-http-checker/0.0.1" "-"`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        200,
			Size:          741,
			Referer:       "-",
			UserAgent:     "mackerel-http-checker/0.0.1",
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] common",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200 741`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        200,
			Size:          741,
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] common (tab delimiter)",
		Input: "10.0.0.11\t-\t-\t[11/Jun/2017:05:56:04 +0900]\t" + `"GET / HTTP/1.1"` + "\t200\t741",
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        200,
			Size:          741,
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] common with empty response",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 204 -`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        204,
			Size:          0,
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] common with vhost",
		Input: `log.example.com 10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 404 741`,
		Output: Log{
			VirtualHost:   "log.example.com",
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        404,
			Size:          741,
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] common with username contains white space",
		Input: `10.0.0.11 - Songmu Yaxing [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200 741`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "Songmu Yaxing",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        200,
			Size:          741,
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] common with virtual host and username contains white space",
		Input: `test.example.com 10.0.0.11 - Songmu Yaxing [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200 741`,
		Output: Log{
			VirtualHost:   "test.example.com",
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "Songmu Yaxing",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        200,
			Size:          741,
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] unescape",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET /?foo=bar HTTP/1.1" 200 741 "\\\thoge" "UserAgent \"fake\""`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET /?foo=bar HTTP/1.1",
			Status:        200,
			Size:          741,
			Referer:       "\\\thoge",
			UserAgent:     `UserAgent "fake"`,
			Method:        "GET",
			RequestURI:    "/?foo=bar",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] unescape(trailing space after escaped double quote)",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET /?foo=bar HTTP/1.1" 200 741 "\" "`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET /?foo=bar HTTP/1.1",
			Status:        200,
			Size:          741,
			Referer:       `" `,
			Method:        "GET",
			RequestURI:    "/?foo=bar",
			Protocol:      "HTTP/1.1",
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
			UserAgent:   "mackerel-agent/0.31.2 (Revision 775fad2)",
			ReqTime:     pfloat64(0.087),
			TakenSec:    pfloat64(0.087),
			Method:      "POST",
			RequestURI:  "/api/v0/tsdb",
			Protocol:    "HTTP/1.1",
		},
	},
	{
		Name: "ltsv with null apptime",
		Input: "time:[08/Mar/2017:14:12:40 +0900]\t" +
			"host:192.0.2.1\t" +
			"req:POST /api/v0/tsdb HTTP/1.1\t" +
			"status:200\t" +
			"size:36\t" +
			"ua:mackerel-agent/0.31.2 (Revision 775fad2)\t" +
			"reqtime:0.087\t" +
			"taken_sec:0.087\t" +
			"apptime:-\t" +
			"vhost:mackerel.io",
		Output: Log{
			VirtualHost: "mackerel.io",
			Host:        "192.0.2.1",
			Time:        time.Date(2017, time.March, 8, 14, 12, 40, 0, loc),
			TimeStr:     "[08/Mar/2017:14:12:40 +0900]",
			Request:     "POST /api/v0/tsdb HTTP/1.1",
			Status:      200,
			Size:        36,
			UserAgent:   "mackerel-agent/0.31.2 (Revision 775fad2)",
			ReqTime:     pfloat64(0.087),
			TakenSec:    pfloat64(0.087),
			Method:      "POST",
			RequestURI:  "/api/v0/tsdb",
			Protocol:    "HTTP/1.1",
		},
	},
	{
		Name: "ltsv filled uri, protocol and method",
		Input: "time:08/Mar/2017:14:12:40 +0900\t" +
			"host:192.0.2.1\t" +
			"req:POST /api/v0/tsdb HTTP/1.1\t" +
			"status:200\t" +
			"size:36\t" +
			"ua:mackerel-agent/0.31.2 (Revision 775fad2)\t" +
			"uri:/api/\t" +
			"protocol:HTTP/1.0\t" +
			"method:GET\t" +
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
			UserAgent:   "mackerel-agent/0.31.2 (Revision 775fad2)",
			ReqTime:     pfloat64(0.087),
			TakenSec:    pfloat64(0.087),
			Method:      "GET",
			RequestURI:  "/api/",
			Protocol:    "HTTP/1.0",
		},
	},
	{
		Name: "LTSV(host field at the beginning)",
		Input: "host:192.0.2.1\t" +
			"time:08/Mar/2017:14:12:40 +0900\t" +
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
			UserAgent:   "mackerel-agent/0.31.2 (Revision 775fad2)",
			ReqTime:     pfloat64(0.087),
			TakenSec:    pfloat64(0.087),
			Method:      "POST",
			RequestURI:  "/api/v0/tsdb",
			Protocol:    "HTTP/1.1",
		},
	},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Logf("testing: %s\n", tt.Name)
			if strings.Contains(tt.Name, "(TODO)") {
				t.Skipf("skip test: %s", tt.Name)
			}
			l, err := Parse(tt.Input)
			if err != nil {
				t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
				return
			}
			if !reflect.DeepEqual(*l, tt.Output) {
				t.Errorf("%s(parse): \n out =%+v\n want %+v", tt.Name, *l, tt.Output)
			}
		})
	}
}

var parsersParseTests = []struct {
	Name   string
	Input  string
	Output Log
}{
	{
		Name:  "[Apache] common",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200 741`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / HTTP/1.1",
			Status:        200,
			Size:          741,
			Method:        "GET",
			RequestURI:    "/",
			Protocol:      "HTTP/1.1",
		},
	},
	{
		Name:  "[Apache] invalid request",
		Input: `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / hoge HTTP/1.1" 200 741`,
		Output: Log{
			Host:          "10.0.0.11",
			RemoteLogname: "-",
			User:          "-",
			Time:          time.Date(2017, time.June, 11, 5, 56, 4, 0, loc),
			Request:       "GET / hoge HTTP/1.1",
			Status:        200,
			Size:          741,
			Method:        "",
			RequestURI:    "",
			Protocol:      "",
		},
	},
	{
		Name: "LTSV",
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
			UserAgent:   "mackerel-agent/0.31.2 (Revision 775fad2)",
			ReqTime:     pfloat64(0.087),
			TakenSec:    pfloat64(0.087),
			Method:      "POST",
			RequestURI:  "/api/v0/tsdb",
			Protocol:    "HTTP/1.1",
		},
	},
	{
		Name: "LTSV (invalid request)",
		Input: "time:08/Mar/2017:14:12:40 +0900\t" +
			"host:192.0.2.1\t" +
			"req:POST /api/v0/tsdb hoge HTTP/1.1\t" +
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
			Request:     "POST /api/v0/tsdb hoge HTTP/1.1",
			Status:      200,
			Size:        36,
			UserAgent:   "mackerel-agent/0.31.2 (Revision 775fad2)",
			ReqTime:     pfloat64(0.087),
			TakenSec:    pfloat64(0.087),
			Method:      "",
			RequestURI:  "",
			Protocol:    "",
		},
	},
}

func TestParsersParse(t *testing.T) {
	ps := Parsers{
		Apache: &Apache{Loose: true},
		LTSV:   &LTSV{Loose: true},
	}
	for _, tt := range parsersParseTests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Logf("testing: %s\n", tt.Name)
			if strings.Contains(tt.Name, "(TODO)") {
				t.Skipf("skip test: %s", tt.Name)
			}
			l, err := ps.Parse(tt.Input)
			if err != nil {
				t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
				return
			}
			if !reflect.DeepEqual(*l, tt.Output) {
				t.Errorf("%s(parse): \n out =%+v\n want %+v", tt.Name, *l, tt.Output)
			}
		})
	}
}
