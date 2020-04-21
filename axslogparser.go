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
	VirtualHost     string `ltsv:"vhost"`
	Host            string
	RemoteLogname   string `ltsv:"ident"`
	User            string
	Time            time.Time `ltsv:"-"`
	TimeStr         string    `ltsv:"time"`
	Request         string    `ltsv:"req"`
	Status          int
	Size            uint64
	Referer         string
	UserAgent       string `ltsv:"ua"`
	ReqTime         *float64
	ReqTimeMicroSec *float64 `ltsv:"reqtime_microsec"`
	AppTime         *float64
	TakenSec        *float64 `ltsv:"taken_sec"` // Hatena specific
	ForwardedFor    string
	RequestURI      string `ltsv:"uri"`
	Protocol        string
	Method          string
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

// Parsers represents a set of available parsers
type Parsers struct {
	Apache *Apache
	LTSV   *LTSV
}

// GuessParser guesses the parser from line
func (ps Parsers) GuessParser(line string) (Parser, *Log, error) {
	if strings.Contains(line, "\thost:") || strings.Contains(line, "\ttime:") {
		l, err := ps.LTSV.Parse(line)
		if err == nil {
			return ps.LTSV, l, err
		}
	}
	l, err := ps.Apache.Parse(line)
	if err != nil {
		return nil, nil, err
	}
	return ps.Apache, l, nil
}

// Parse log line
func (ps Parsers) Parse(line string) (*Log, error) {
	_, l, err := ps.GuessParser(line)
	return l, err
}

// GuessParser guesses the parser from line, uses default parsers for each format
func GuessParser(line string) (Parser, *Log, error) {
	ps := Parsers{
		LTSV:   &LTSV{},
		Apache: &Apache{},
	}
	return ps.GuessParser(line)
}

// Parse log line using default parsers
func Parse(line string) (*Log, error) {
	_, l, err := GuessParser(line)
	return l, err
}
