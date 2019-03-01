package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	e *Entry
	http.Handler
}

func NewProxy(e *Entry) (*Proxy, error) {
	proxy := &Proxy{e: e}
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

func (p *Proxy) ServeConnect(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("CONNECT: ignoring")
}
