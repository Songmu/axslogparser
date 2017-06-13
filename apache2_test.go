package axslogparser

import (
	"testing"
)

func BenchmarkApache2_Parse(b *testing.B) {
	p := &Apache2{}
	for i := 0; i < b.N; i++ {
		for _, v := range benchInputs {
			p.Parse(v)
		}
	}
}
