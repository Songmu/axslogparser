package axslogparser

import (
	"fmt"
	"strings"
	"time"
)

// Parser is the interface for accesslog
type Parser interface {
	Parse(string) (*Log, error)
}

// Log is the struct stored parsed result of single line of accesslog
type Log struct {
	VirtualHost   string `ltsv:"vhost"`
	Host          string
	RemoteLogname string `ltsv:"ident"`
	User          string
	Time          time.Time `ltsv:"-"`
	TimeStr       string    `ltsv:"time"`
	Request       string    `ltsv:"req"`
	Status        int
	Size          uint64
	Referer       string
	UserAgent     string `ltsv:"ua"`
	ReqTime       *float64
	AppTime       *float64
	TakenSec      *float64 `ltsv:"taken_sec"` // Hatena specific
	ForwardedFor  string
	RequestURI    string `ltsv:"uri"`
	Protocol      string
	Method        string
}

func (l *Log) breakdownRequest() error {
	if l.RequestURI != "" && l.Protocol != "" && l.Method != "" {
		return nil
	}
	stuff := strings.Fields(l.Request)
	if len(stuff) != 3 {
		return fmt.Errorf("invalid request: %s", l.Request)
	}
	if len(stuff) > 0 && l.Method == "" {
		l.Method = stuff[0]
	}
	if len(stuff) > 1 && l.RequestURI == "" {
		l.RequestURI = stuff[1]
	}
	if len(stuff) > 2 && l.Protocol == "" {
		l.Protocol = stuff[2]
	}
	return nil
}

const clfTimeLayout = "02/Jan/2006:15:04:05 -0700"

// GuessParser guesses the parser from line
func GuessParser(line string) (Parser, *Log, error) {
	var p Parser
	if strings.Contains(line, "\thost:") || strings.Contains(line, "\ttime:") {
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

// Parse log line
func Parse(line string) (*Log, error) {
	_, l, err := GuessParser(line)
	return l, err
}
