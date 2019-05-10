package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/octavore/naga/service"
)

const (
	defaultConfigDir = ".config"
	version          = "1.1.1"
)

func main() {
	service.Run(&App{})
}

type App struct {
	config       *Config
	handlers     []*Proxy
	handlerIndex map[string]int // use for colors
}

func (a *App) Init(c *service.Config) {
	a.handlerIndex = make(map[string]int)
	a.addCommands(c)
	c.SetDefaultCommand("start")

	c.Start = func() {
		err := a.ensureConfig()
		if err != nil {
			log.Fatalln(err)
		}

		err = a.loadConfig()
		if err != nil {
			log.Fatalln(err)
		}
		if a.config.Addr == "" {
			a.config.Addr = "localhost:7999"
		}
		if a.config.TLSAddr == "" {
			a.config.TLSAddr = "localhost:7998"
		}
		if len(a.config.Entries) == 0 {
			fmt.Println("no entries found in config.json; exiting")
			os.Exit(1)
		}
		for i, e := range a.config.Entries {
			p, err := NewProxy(e)
			if err != nil {
				log.Fatalln(err)
			}
			a.handlers = append(a.handlers, p)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			a.handlerIndex[e.Source] = i
			i := a.handlerIndex[e.Source] % len(colors)
			fmt.Printf("loaded: %s => %s\n", colors[i](e.Source), e.dest())
		}

		fmt.Printf("proxy URL: http://%s/proxy.pac\n", a.config.Addr)
		a.startServer()
	}
}

func (a *App) startServer() {
	err := a.startTLSProxy()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = http.ListenAndServe(a.config.Addr, a)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
