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
