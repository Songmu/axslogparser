package axslogparser

import (
	"strings"
	"time"
)

type Parser interface {
	Parse(string) (*Log, error)
}

type Log struct {
	VirtualHost  string `ltsv:"vhost"`
	Host         string
	User         string
	Time         time.Time `ltsv:"-"`
	TimeStr      string    `ltsv:"time"`
	Request      string    `ltsv:"req"`
	Status       int
	Size         uint64
	Referer      string
	UA           string
	ReqTime      *float64
	AppTime      *float64
	TakenSec     *float64 `ltsv:"taken_sec"` // Hatena specific
	ForwardedFor string
	URI          string
	Protocol     string
	Method       string
}

func (l *Log) breakdownRequest() {
	if l.URI != "" && l.Protocol != "" && l.Method != "" {
		return
	}
	stuff := strings.Fields(l.Request)
	if len(stuff) > 0 && l.Method == "" {
		l.Method = stuff[0]
	}
	if len(stuff) > 1 && l.URI == "" {
		l.URI = stuff[1]
	}
	if len(stuff) > 2 && l.Protocol == "" {
		l.Protocol = stuff[2]
	}
	return
}

const clfTimeLayout = "02/Jan/2006:15:04:05 -0700"

func GuessParser(line string) (Parser, *Log, error) {
	var p Parser
	if strings.Contains(line, "\thost:") {
		p = &LTSV{}
		l, err := p.Parse(line)
		if err == nil {
			return p, l, err
		}
	}
	p = &Apache{}
	l, err := p.Parse(line)
	if err != nil {
		return nil, nil, err
	}
	return p, l, nil
}

func Parse(line string) (*Log, error) {
	_, l, err := GuessParser(line)
	return l, err
}
