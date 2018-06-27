package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Entry represents a host to route to either another host, or a
// file folder on disk.
type Entry struct {
	Source     string `json:"host"`
	DestHost   string `json:"dest,omitempty"`
	DestFolder string `json:"dest_folder,omitempty"`
}

func (e *Entry) handle() (http.Handler, error) {
	if e.DestHost != "" {
		if !strings.HasPrefix(e.DestHost, "http") {
			e.DestHost = "http://" + e.DestHost
		}
		u, err := url.Parse(e.DestHost)
		if err != nil {
			panic(err)
		}
		return httputil.NewSingleHostReverseProxy(u), nil
	}
	if e.DestFolder != "" {
		return http.FileServer(http.Dir(e.DestFolder)), nil
	}
	return nil, fmt.Errorf("entry: no dest or dest_folder provided")
}

func (e *Entry) dest() string {
	if e.DestHost != "" {
		return e.DestHost
	}
	return e.DestFolder
}
