package main

import (
	"log"
	"net/http"

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

// ServeHTTP checks the host of the request and calls the registered handler. If the host
// is not registered, a 404 error is returned
func (a *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if h := a.handlers[req.Host]; h != nil {
		i := a.handlerIndex[req.Host] % len(colors)
		pth := req.URL.Path
		if req.URL.RawQuery != "" {
			pth += "?" + req.URL.RawQuery
		}
		log.Println(colors[i](req.Host), pth)
		h.ServeHTTP(rw, req)
		return
	}
	log.Println(color.RedString(req.Host), req.URL.String())
	http.Error(rw, "not mapped: "+req.Host, http.StatusNotFound)
}
