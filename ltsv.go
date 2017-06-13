package axslogparser

import (
	"time"

	"github.com/Songmu/go-ltsv"
	"github.com/pkg/errors"
)

type LTSV struct {
}

func (lv *LTSV) Parse(line string) (*Log, error) {
	l := &Log{}
	err := ltsv.Unmarshal([]byte(line), l)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ltsvlog")
	}
	l.Time, _ = time.Parse(clfTimeLayout, l.TimeStr)
	l.breakdownRequest()
	return l, nil
}
