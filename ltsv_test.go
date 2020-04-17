package axslogparser

import (
	"strings"
	"testing"
)

var parseLTSVErrorTests = []struct {
	Name           string
	Input          string
	ContainsString string
	HardError      bool
}{
	{
		Name:           "broken",
		Input:          "hoge",
		ContainsString: "(not a ltsv)",
		HardError:      true,
	},
	{
		Name: "invalid request",
		Input: "time:08/Mar/2017:14:12:40 +0900\t" +
			"host:192.0.2.1\t" +
			"req:POST /api/v0/tsdb hoge HTTP/1.1\t" +
			"status:200\t" +
			"size:36\t" +
			"ua:mackerel-agent/0.31.2 (Revision 775fad2)\t" +
			"reqtime:0.087\t" +
			"taken_sec:0.087\t" +
			"vhost:mackerel.io",
		ContainsString: "(invalid request)",
	},
	{
		Name: "unmarshal error",
		Input: "time:08/Mar/2017:14:12:40 +0900\t" +
			"host:192.0.2.1\t" +
			"req:POST /api/v0/tsdb HTTP/1.1\t" +
			"status:200\t" +
			"size:36\t" +
			"ua:mackerel-agent/0.31.2 (Revision 775fad2)\t" +
			"reqtime:0.087\t" +
			"taken_sec:0.087\t" +
			"vhost:mackerel.io\t" +
			"apptime:0.021, 0.042", // not a single number
		ContainsString: "ltsv: cannot unmarshal",
	},
}

func TestLTSV_ParseError(t *testing.T) {
	psr := &LTSV{}
	for _, tt := range parseLTSVErrorTests {
		t.Logf("testing: %s", tt.Name)
		if _, err := psr.Parse(tt.Input); err == nil {
			t.Errorf("%s: error should be occured but nil", tt.Name)
		} else if !strings.Contains(err.Error(), tt.ContainsString) {
			t.Errorf("%s: error should be contained %q, but: %s", tt.Name, tt.ContainsString, err)
		}
	}
}

func TestLTSV_ParseError_loose(t *testing.T) {
	psr := &LTSV{Loose: true}
	for _, tt := range parseLTSVErrorTests {
		t.Logf("testing: %s", tt.Name)
		if tt.HardError {
			if _, err := psr.Parse(tt.Input); err == nil {
				if tt.HardError {
					t.Errorf("%s: error should be occured but nil", tt.Name)
				}
			} else if !strings.Contains(err.Error(), tt.ContainsString) {
				t.Errorf("%s: error should be contained %q, but: %s", tt.Name, tt.ContainsString, err)
			}
		} else {
			if _, err := psr.Parse(tt.Input); err != nil {
				if err != nil {
					t.Errorf("%s: error should not be occured but: %s", tt.Name, err.Error())
				}
			}
		}
	}
}
