package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	if p.e.DestHost == "" {
		hj, isHJ := rw.(http.Hijacker)
		if !isHJ {
			fmt.Println("CONNECT error: cannot hijack CONNECT")
			return
		}

		// hijack connection
		c, br, err := hj.Hijack()
		if err != nil {
			fmt.Println("CONNECT error", err)
			http.Error(rw, err.Error(), 500)
			return
		}
		defer c.Close()

		// connect to backend
		be, err := net.Dial("tcp", p.e.DestHost)
		if err != nil {
			fmt.Println("CONNECT error:", err)
			return
		}
		defer be.Close()

		// write request to backend
		if err := req.Write(be); err != nil {
			log.Printf("websocket backend write request: %v", err)
			http.Error(rw, err.Error(), 500)
			return
		}

		// connect the two pipes
		errc := make(chan error, 1)
		go func() {
			n, err := io.Copy(be, br) // backend <- buffered reader
			if err != nil {
				err = fmt.Errorf("websocket: to copy backend from buffered reader: %v, %v", n, err)
			}
			errc <- err
		}()
		go func() {
			n, err := io.Copy(c, be) // raw conn <- backend
			if err != nil {
				err = fmt.Errorf("websocket: to raw conn from backend: %v, %v", n, err)
			}
			errc <- err
		}()
		if err := <-errc; err != nil {
			log.Print(err)
		}
		return
	}
}
