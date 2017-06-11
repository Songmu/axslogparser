package axslogparser

import (
	"reflect"
	"testing"
	"time"
)

var loc = func() *time.Location {
	l, _ := time.LoadLocation("Asia/Tokyo")
	return l
}()

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
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		l, err := Parse(tt.Input, loc)
		if err != nil {
			t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
			continue
		}
		if !reflect.DeepEqual(*l, tt.Output) {
			t.Errorf("%s(parse): \n out =%+v\n want %+v", tt.Name, *l, tt.Output)
		}
	}
}
