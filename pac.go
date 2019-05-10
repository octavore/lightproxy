package main

import (
	"log"
	"net/http"
	"net/url"
)

func (a *App) pacFile(rw http.ResponseWriter, req *http.Request) {
	u, err := url.Parse(a.config.Addr)
	if err != nil {
		log.Fatalln("error parsing addr in config file:", err)
	}
	rw.Write([]byte(`function FindProxyForURL(url, host) {
	if (shExpMatch(host, "*.` + a.config.TLD + `")) {
		return "PROXY 127.0.0.1:` + u.Port() + `";
	}
	return "DIRECT";
}
`))
}
