package main

import (
	"io"
	"net"
	"net/http"
	"time"
)

// serveConnect handles a CONNECT call on a.config.Addr and connects it to the
// tlsProxy on a.config.TLSAddr
func (a *App) serveConnect(rw http.ResponseWriter, req *http.Request) {
	destConn, err := net.DialTimeout("tcp", a.config.TLSAddr, 10*time.Second)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}
	rw.WriteHeader(http.StatusOK)
	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
