package axslogparser

import (
	"strings"
	"time"
)

type Parser interface {
	Parse(string) (*Log, error)
}

type Log struct {
	VirtualHost string
	Host        string
	User        string
	Time        time.Time
	Request     string
	Status      int
	Size        uint64
	Referer     string
	UA          string
}

const clfTimeLayout = "02/Jan/2006:15:04:05 -0700"

func GuessParser(line string) (Parser, *Log, error) {
	var p Parser
	if strings.Contains(line, "\thost:") {
		p = &LTSV{}
		l, err := p.Parse(line)
		if err != nil {
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
