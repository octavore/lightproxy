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
	version          = "1.1.0"
)

func main() {
	service.Run(&App{})
}

type App struct {
	config       *Config
	handlers     map[string]*Proxy
	handlerIndex map[string]int // use for colors
}

func (a *App) Init(c *service.Config) {
	a.handlers = make(map[string]*Proxy)
	a.handlerIndex = make(map[string]int)
	a.addCommands(c)
	c.SetDefaultCommand("start")

	c.Start = func() {
		err := a.loadConfig()
		if err != nil {
			log.Fatalln(err)
		}
		if len(a.config.Entries) == 0 {
			fmt.Println("no entries found in config.json; exiting")
			os.Exit(1)
		}
		for i, e := range a.config.Entries {
			a.handlers[e.Source], err = NewProxy(e)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			a.handlerIndex[e.Source] = i
			i := a.handlerIndex[e.Source] % len(colors)
			fmt.Printf("loaded: %s => %s\n", colors[i](e.Source), e.dest())
		}

		fmt.Printf("proxy URL: http://%s/proxy.pac\n", a.config.Addr)
		err = http.ListenAndServe(a.config.Addr, a)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
