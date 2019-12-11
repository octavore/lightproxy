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

var defaultConfig = &Config{
	Addr:    "localhost:7999",
	TLSAddr: "localhost:7998",
	TLD:     "wip",
	Entries: []*Entry{{
		Source:   "example.wip",
		DestHost: "localhost:8000",
	}},
}

type configManager struct {
	searchPaths []string
}

// newConfigManager sets up all the search paths and configPath
func newConfigManager() (*configManager, error) {
	cm := &configManager{}
	c := os.Getenv("XDG_CONFIG_HOME")
	if c != "" {
		cm.searchPaths = append(cm.searchPaths, path.Join(c, "lightproxy"))
	}
	cm.searchPaths = append(cm.searchPaths, getHomeConfigDir())
	return cm, nil
}

// configPath() returns the active config file path, config file dir
// and whether it exists or not.
// This checks all search paths for an existing config file
// XDG_CONFIG_HOME is the preferred path, but also fallback
// gracefully to $HOME/.config
func (cm *configManager) configPath() (string, string, bool) {
	// default config path is the first search path
	configDir := cm.searchPaths[0]

	for _, dir := range cm.searchPaths {
		configPath := path.Join(dir, "config.json")
		fi, err := os.Stat(configPath)
		if fi != nil && err == nil {
			return configPath, dir, true
		}
		if !os.IsNotExist(err) {
			// fmt.Printf("unknown error: %s\n", err)
		}
	}
	return path.Join(configDir, "config.json"), configDir, false
}

func (cm *configManager) ensureAndLoad() (*Config, error) {
	err := cm.ensure()
	if err != nil {
		return nil, err
	}
	configPath, _, exists := cm.configPath()
	config := &Config{}
	if exists {
		f, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(f, config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

// writeConfig writes to the existing config file, or the first search path
func (cm *configManager) writeConfig(config *Config) error {
	configPath, _, _ := cm.configPath()
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, b, os.ModePerm)
}

func (cm *configManager) ensure() error {
	configPath, configDir, exists := cm.configPath()
	if exists {
		return nil
	} else {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create dir %s: %s", configDir, err)
		}
	}
	err := cm.writeConfig(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to to create config.json file: %s", err)
	}
	fmt.Printf("created config.json file: %s\n", configPath)
	return nil
}

func getHomeConfigDir() string {
	u, err := user.Current()
	if uid := os.Getenv("SUDO_UID"); uid != "" {
		u, err = user.LookupId(uid)
	}
	if err != nil {
		panic(err)
	}
	return path.Join(u.HomeDir, ".config", "lightproxy")
}
