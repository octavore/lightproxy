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
	if h := a.handlers[host]; h != nil {
		i := a.handlerIndex[host] % len(colors)
		pth := req.URL.Path
		if req.URL.RawQuery != "" {
			pth += "?" + req.URL.RawQuery
		}
		log.Println(colors[i](host), pth)

		if req.Method == "CONNECT" {
			h.ServeConnect(rw, req)
		} else {
			h.ServeHTTP(rw, req)
		}
		return
	}
	log.Println(color.RedString(host), req.URL.String())
	http.Error(rw, "not mapped: "+host, http.StatusNotFound)
}
