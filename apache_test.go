package axslogparser

import (
	"strings"
	"testing"
)

var benchInputs []string

func init() {
	for _, v := range parseTests {
		if strings.HasPrefix(v.Name, "[Apache]") {
			benchInputs = append(benchInputs, v.Input)
		}
	}
}

func BenchmarkApache_Parse(b *testing.B) {
	p := &Apache{}
	for i := 0; i < b.N; i++ {
		for _, v := range benchInputs {
			p.Parse(v)
		}
	}
}

var parseErrorTests = []struct {
	Name           string
	Input          string
	ContainsString string
}{
	{
		Name:           "broken",
		Input:          "hoge",
		ContainsString: "(not matched)",
	},
	{
		Name:           "invalid request",
		Input:          `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / hoge HTTP/1.1" 200 741 "-" "mackerel-http-checker/0.0.1"`,
		ContainsString: "(invalid request)",
	},
	{
		Name:           "invalid request",
		Input:          `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 200`,
		ContainsString: "(invalid status or size)",
	},
	{
		Name:           "invalid request",
		Input:          `10.0.0.11 - - [11/Jun/2017:05:56:04 +0900] "GET / HTTP/1.1" 2xx 741 "-" "mackerel-http-checker/0.0.1"`,
		ContainsString: "(invalid status:",
	},
}

func TestParse_error(t *testing.T) {
	psr := &Apache{}
	for _, tt := range parseErrorTests {
		t.Logf("testing: %s", tt.Name)
		if _, err := psr.Parse(tt.Input); err == nil {
			t.Errorf("%s: error should be occured but nil", tt.Name)
		} else if !strings.Contains(err.Error(), tt.ContainsString) {
			t.Errorf("%s: error should be contained %q, but: %s", tt.Name, tt.ContainsString, err)
		}
	}
}
