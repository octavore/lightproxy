package main

import "net/http"

func (a *App) pacFile(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(`function FindProxyForURL(url, host) {
	if (shExpMatch(host, "*.` + a.config.TLD + `")) {
		return "PROXY 127.0.0.1:7999";
	}
	return "DIRECT";
}
`))
}
