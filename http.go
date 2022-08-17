package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type formatter func(format string, a ...interface{}) string

var colors = []formatter{
	color.YellowString,
	color.GreenString,
	color.MagentaString,
	color.CyanString,
	color.BlueString,
}

var portRegexp = regexp.MustCompile(`:\d+$`)

// ServeHTTP checks the host of the request and calls the registered handler. If the host
// is not registered, a 404 error is returned
func (a *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/proxy.pac" {
		a.pacFile(rw, req)
		return
	}

	host := req.Host
	host = strings.TrimPrefix(host, "//")
	host = portRegexp.ReplaceAllString(host, "")
	for i, h := range a.handlers {
		if h.Match(req) {
			c := i % len(colors)
			pth := req.URL.Path
			if req.URL.RawQuery != "" {
				pth += "?" + req.URL.RawQuery
			}

			if req.Method == "CONNECT" {
				log.Printf("%s (received CONNECT)", colors[c](host))
				a.serveConnect(rw, req, h.e)
			} else {
				log.Println(colors[c](host), pth)
				h.ServeHTTP(rw, req)
			}
			return
		}
	}

	log.Println(color.RedString(host), req.URL.String())
	http.Error(rw, "not mapped: "+host, http.StatusNotFound)
}
