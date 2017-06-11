package axslogparser

import "fmt"

type LTSV struct {
}

func (lv *LTSV) Parse(line string) (*Log, error) {
	return nil, fmt.Errorf("not implemented")
}
