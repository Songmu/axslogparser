package axslogparser

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Apache log parser
type Apache struct {
}

var logRe = regexp.MustCompile(
	`^(?:(\S+(?:,\s\S+)*)\s)?` + // %v(The canonical ServerName/virtual host) - 192.168.0.1 or 192.168.0.1,192.168.0.2, 192.168.0.3
		`(\S+)\s` + // %h(Remote Hostname) $remote_addr
		`(\S+)\s` + // %l(Remote Logname)
		`([\S\s]+)\s` + // $remote_user
		`\[(\d{2}/\w{3}/\d{2}(?:\d{2}:){3}\d{2} [-+]\d{4})\]\s` + // $time_local
		`(.*)`)

// Parse for Parser interface
func (ap *Apache) Parse(line string) (*Log, error) {
	matches := logRe.FindStringSubmatch(line)
	if len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse apachelog (not matched): %s", line)
	}
	l := &Log{
		VirtualHost:   matches[1],
		Host:          matches[2],
		RemoteLogname: matches[3],
		User:          matches[4],
	}
	if l.Host == "-" && l.VirtualHost != "" {
		l.Host = l.VirtualHost
		l.VirtualHost = ""
		l.User = fmt.Sprintf("%s %s", l.RemoteLogname, l.User)
		l.RemoteLogname = "-"
	}

	l.Time, _ = time.Parse(clfTimeLayout, matches[5])
	var rest string

	l.Request, rest = takeQuoted(matches[6])
	if err := l.breakdownRequest(); err != nil {
		return nil, errors.Wrapf(err, "failed to parse apachelog (invalid request): %s", line)
	}
	matches = strings.Fields(rest)
	if len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse apachelog (invalid status or size): %s", line)
	}
	l.Status, _ = strconv.Atoi(matches[0])
	if l.Status < 100 || 600 <= l.Status {
		return nil, fmt.Errorf("failed to parse apachelog (invalid status: %s): %s", matches[0], line)
	}
	l.Size, _ = strconv.ParseUint(matches[1], 10, 64)
	l.Referer, rest = takeQuoted(rest)
	l.UserAgent, _ = takeQuoted(rest)
	return l, nil
}

func takeQuoted(line string) (string, string) {
	if line == "" {
		return "", ""
	}
	i := 0
	for ; i < len(line); i++ {
		if line[i] == '"' {
			i++
			break
		}
	}
	if i == len(line) {
		return "", ""
	}
	buf := &bytes.Buffer{}
	escaped := false
	for ; i < len(line); i++ {
		c := line[i]
		if !escaped {
			if c == '"' {
				break
			}
			if c == '\\' {
				escaped = true
				continue
			}
			buf.WriteByte(c)
			continue
		}
		escaped = false
		switch c {
		case 'n':
			buf.WriteByte('\n')
		case 't':
			buf.WriteByte('\t')
		case '\\':
			buf.WriteByte('\\')
		case '"':
			buf.WriteByte('"')
		default:
			buf.WriteByte(c)
		}
	}
	return buf.String(), line[i+1:]
}
