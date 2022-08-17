package main

import (
	"io"
	"net"
	"net/http"
	"time"
	"strings"
)

// serveConnect handles a CONNECT call on a.config.Addr and connects it to the
// tlsProxy on a.config.TLSAddr
func (a *App) serveConnect(rw http.ResponseWriter, req *http.Request, entry *Entry) {
	destHost := strings.TrimPrefix(entry.DestHost, "http://")
	if req.URL.Port() == "443" {
		destHost = a.config.TLSAddr
	}

	destConn, err := net.DialTimeout("tcp", destHost, 10*time.Second)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}

	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
	}

	clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	_, _ = io.Copy(destination, source)
}
