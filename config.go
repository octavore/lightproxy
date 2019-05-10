package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type Config struct {
	TLD       string   `json:"tld"`
	Addr      string   `json:"addr"`
	TLSAddr   string   `json:"tls_addr"`
	CAKeyFile string   `json:"ca_key_file"`
	Entries   []*Entry `json:"entries"`
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

func (a *App) writeConfig() error {
	b, err := json.MarshalIndent(a.config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(a.configPath(), b, os.ModePerm)
}
