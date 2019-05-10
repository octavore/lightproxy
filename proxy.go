package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

type Proxy struct {
	e *Entry
	http.Handler
	re *regexp.Regexp
}

func NewProxy(e *Entry) (*Proxy, error) {
	proxy := &Proxy{e: e}

	// quote all characters in r and convert * to non empty segments
	r := regexp.QuoteMeta(e.Source)
	r = "^" + strings.ReplaceAll(r, `\*`, `[^.]+`) + "$"
	proxy.re = regexp.MustCompile(r)

	if e.DestHost != "" {
		if !strings.HasPrefix(e.DestHost, "http") {
			e.DestHost = "http://" + e.DestHost
		}
		target, err := url.Parse(e.DestHost)
		if err != nil {
			panic(err)
		}
		proxy.Handler = httputil.NewSingleHostReverseProxy(target)
	} else if e.DestFolder != "" {
		proxy.Handler = http.FileServer(http.Dir(e.DestFolder))
	} else {
		return nil, fmt.Errorf("entry: no dest or dest_folder provided")
	}

	return proxy, nil
}

func (p *Proxy) Match(req *http.Request) bool {
	host := req.Host
	host = strings.TrimPrefix(host, "//")
	host = portRegexp.ReplaceAllString(host, "")
	return p.re.MatchString(host)
}
