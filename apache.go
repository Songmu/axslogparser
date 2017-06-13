package axslogparser

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Apache log parser
type Apache struct {
}

var logRe = regexp.MustCompile(
	`(?:(?P<vhost>\S+)\s)?` + // %v(The canonical ServerName/virtual host)
		`(?P<remote_addr>\S+)\s` + // %h(Remote Hostname) $remote_addr
		`\S+\s+` + // %l(Remote Logname)
		`(?P<remote_user>\S+)\s` + // $remote_user
		`\[(?P<time_local>\d{2}/\w{3}/\d{2}(?:\d{2}:){3}\d{2} [-+]\d{4})\]\s` + // $time_local
		`(?P<rest>.*)`)

// Parse for Parser interface
func (ap *Apache) Parse(line string) (*Log, error) {
	matches := logRe.FindStringSubmatch(line)
	if len(matches) < 1 {
		return nil, fmt.Errorf("not matched")
	}
	var rest string
	l := &Log{}
	for i, name := range logRe.SubexpNames() {
		switch name {
		case "vhost":
			l.VirtualHost = matches[i]
		case "remote_addr":
			l.Host = matches[i]
		case "remote_user":
			l.User = matches[i]
		case "time_local":
			l.Time, _ = time.Parse(clfTimeLayout, matches[i])
		case "rest":
			rest = matches[i]
		}
	}
	l.Request, rest = takeQuoted(rest)
	matches = strings.Fields(rest)
	if len(matches) > 1 {
		l.Status, _ = strconv.Atoi(matches[0])
		l.Size, _ = strconv.ParseUint(matches[1], 10, 64)
	}
	l.Referer, rest = takeQuoted(rest)
	l.UA, _ = takeQuoted(rest)
	l.breakdownRequest()
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
