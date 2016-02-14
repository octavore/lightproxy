package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"

	"github.com/octavore/naga/service"
)

const (
	defaultConfigDir = ".config"
	version          = "0.1.0"
)

func main() {
	service.Run(&App{})
}

type App struct {
	config       *Config
	handlers     map[string]http.Handler
	handlerIndex map[string]int
}

func (a *App) Init(c *service.Config) {
	a.handlers = make(map[string]http.Handler)
	a.handlerIndex = make(map[string]int)
	c.AddCommand(&service.Command{
		Keyword:    "init",
		ShortUsage: "initialize the config file",
		Usage:      "Initialize a default config file if it doesn't already exist, and print its location",
		Run:        a.cmdInitConfig,
	})

	c.AddCommand(&service.Command{
		Keyword:    "config",
		ShortUsage: "prints the config file",
		Usage:      "Prints the config files",
		Run:        a.cmdPrintConfig,
	})

	c.AddCommand(&service.Command{
		Keyword:    "set-dest <host> <dest>",
		ShortUsage: "map <host> to <dest>",
		Usage:      "Map <host> to <dest>",
		Run:        a.cmdSetHost,
	})

	c.AddCommand(&service.Command{
		Keyword:    "version",
		ShortUsage: "print version",
		Usage:      "Print version",
		Run: func(*service.CommandContext) {
			fmt.Println("lightproxy", version)
		},
	})

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
			a.handlers[e.Source], err = e.handle()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			a.handlerIndex[e.Source] = i
			fmt.Printf("loaded: %s => %s\n", e.Source, e.dest())
		}
		err = http.ListenAndServe(a.config.Addr, a)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func (a *App) configDir() string {
	c := os.Getenv("XDG_CONFIG_HOME")
	if c == "" {
		u, err := user.Current()
		if uid := os.Getenv("SUDO_UID"); uid != "" {
			u, err = user.LookupId(uid)
		}
		if err != nil {
			panic(err)
		}

		c = path.Join(u.HomeDir, defaultConfigDir)
	}
	return path.Join(c, "lightproxy")
}

func (a *App) configPath() string {
	return path.Join(a.configDir(), "config.json")

}

func (a *App) loadConfig() error {
	f, err := ioutil.ReadFile(a.configPath())
	if err != nil {
		return err
	}
	a.config = &Config{}
	return json.Unmarshal(f, a.config)
}
