package axslogparser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Apache2 log parser
type Apache2 struct {
}

var part2 = `"(?P<%s>(?:[^"]*(?:\\")?)*)"(?:\s|$)`

var logRe2 = regexp.MustCompile(
	`(?:(?P<vhost>\S+)\s)?` + // %v(The canonical ServerName/virtual host)
		`(?P<remote_addr>\S+)\s` + // %h(Remote Hostname) $remote_addr
		`\S+\s+` + // %l(Remote Logname)
		`(?P<remote_user>\S+)\s` + // $remote_user
		`\[(?P<time_local>\d{2}/\w{3}/\d{2}(?:\d{2}:){3}\d{2} [-+]\d{4})\]\s` + // $time_local
		fmt.Sprintf(part2, "request") + // $request
		`(?P<status>[0-9]{3})\s` + // $status
		`(?P<body_bytes_sent>-|(?:[0-9]+))(?:$|\s)` + // $body_bytes_sent
		`(?:` + // combined option start
		fmt.Sprintf(part2, "http_referer") + // $http_referer
		fmt.Sprintf(part2, "http_user_agent") + // $http_user_agent
		`)?`) // combined option end

// Parse for Parser interface
func (ap *Apache2) Parse(line string) (*Log, error) {
	matches := logRe2.FindStringSubmatch(line)
	if len(matches) < 1 {
		return nil, fmt.Errorf("not matched")
	}
	l := &Log{}
	for i, name := range logRe2.SubexpNames() {
		switch name {
		case "vhost":
			l.VirtualHost = matches[i]
		case "remote_addr":
			l.Host = matches[i]
		case "remote_user":
			l.User = matches[i]
		case "time_local":
			l.Time, _ = time.Parse(clfTimeLayout, matches[i])
		case "request":
			l.Request = unescape(matches[i])
		case "status":
			l.Status, _ = strconv.Atoi(matches[i])
		case "body_bytes_sent":
			v := matches[i]
			if v == "-" {
				v = "0"
			}
			l.Size, _ = strconv.ParseUint(matches[i], 10, 64)
		case "http_referer":
			l.Referer = unescape(matches[i])
		case "http_user_agent":
			l.UA = unescape(matches[i])
		}
	}
	l.breakdownRequest()
	return l, nil
}
