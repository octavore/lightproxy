package main

import (
	"encoding/json"
	"fmt"
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

func (a *App) ensureConfig() error {
	fi, err := os.Stat(a.configPath())
	if fi != nil && err == nil {
		fmt.Printf("found config file: %s\n", a.configPath())
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("unknown error: %s", err)
	}

	err = os.MkdirAll(a.configDir(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create dir %s: %s", a.configDir(), err)
	}

	f, err := os.Create(a.configPath())
	defer f.Close()
	if err != nil {
		return fmt.Errorf("failed to create dir %s: %s", a.configDir(), err)
	}

	b, err := json.MarshalIndent(&Config{
		Addr:    "localhost:7999",
		TLSAddr: "localhost:7998",
		TLD:     "wip",
		Entries: []*Entry{{
			Source:   "example.wip",
			DestHost: "localhost:8000",
		}},
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to to create config.json file: %s", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("failed to to create config.json file: %s", err)
	}
	fmt.Printf("created config.json file: %s\n", a.configPath())
	return nil
}
