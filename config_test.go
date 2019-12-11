package main

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureConfig(t *testing.T) {
	dirA, err := ioutil.TempDir("", "lightproxy")
	if err != nil {
		t.Fatal(err)
	}

	dirB, err := ioutil.TempDir("", "lightproxy")
	if err != nil {
		t.Fatal(err)
	}

	cm := &configManager{searchPaths: []string{dirA, dirB}}
	configPath, configDir, exists := cm.configPath()
	// should not exist yet
	assert.False(t, exists, "config file should not exist")
	assert.Equal(t, dirA, configDir)
	assert.Equal(t, path.Join(dirA, "config.json"), configPath)

	// ensureAndLoad creates a default file
	config, err := cm.ensureAndLoad()
	assert.NoError(t, err)
	assert.Equal(t, defaultConfig, config)

	// should exist after we called ensure
	configPath, configDir, exists = cm.configPath()
	assert.True(t, exists, "config file should exist")
	assert.Equal(t, dirA, configDir)
	assert.Equal(t, path.Join(dirA, "config.json"), configPath)

	// should work by falling back correctly
	reverseSearchPaths := &configManager{searchPaths: []string{dirB, dirA}}
	configPath, configDir, exists = reverseSearchPaths.configPath()
	assert.True(t, exists, "config file should not exist")
	assert.Equal(t, dirA, configDir)
	assert.Equal(t, path.Join(dirA, "config.json"), configPath)
}
