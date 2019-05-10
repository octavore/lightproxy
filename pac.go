package main

import (
	"net/http"
	"strings"
)

func (a *App) pacFile(rw http.ResponseWriter, req *http.Request) {
	parts := strings.Split(a.config.Addr, ":")
	port := parts[len(parts)-1]
	rw.Write([]byte(`function FindProxyForURL(url, host) {
	if (shExpMatch(host, "*.` + a.config.TLD + `")) {
		return "PROXY 127.0.0.1:` + port + `";
	}
	return "DIRECT";
}
`))
}
