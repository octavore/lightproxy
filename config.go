package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Entry struct {
	Source     string  `json:"host"`
	DestHost   *string `json:"dest,omitempty"`
	DestFolder *string `json:"dest_folder,omitempty"`
}

type Config struct {
	Addr      string  `json:"addr"`
	AdminAddr *string `json:"admin_addr"`
	Entries   []Entry `json:"entries"`
}

func (e *Entry) handle() (http.Handler, error) {
	if e.DestHost != nil {
		if !strings.HasPrefix(*e.DestHost, "http") {
			*e.DestHost = "http://" + *e.DestHost
		}
		u, err := url.Parse(*e.DestHost)
		if err != nil {
			panic(err)
		}
		return httputil.NewSingleHostReverseProxy(u), nil
	}
	if e.DestFolder != nil {
		return http.FileServer(http.Dir(*e.DestFolder)), nil
	}
	return nil, fmt.Errorf("entry: no dest or dest_folder provided")
}

func (e *Entry) dest() string {
	if e.DestHost != nil {
		return *e.DestHost
	}
	if e.DestFolder != nil {
		return *e.DestFolder
	}
	return ""
}
