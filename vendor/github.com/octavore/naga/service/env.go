package service

import (
	"os"
	"strings"
)

// Environment can be configured with the NAGA_ENV environment variable.
// This defaults to EnvUnknown if not otherwise set.s
type Environment int

// EnvVarName is the name of environment variable which contains the app
// environment. Defaults to NAGA_ENV.
var EnvVarName = "NAGA_ENV"

// Environment constants
const (
	EnvProduction Environment = iota
	EnvDevelopment
	EnvTest
	EnvUnknown
)

// EnvMap gives mappings of strings to specific environments.
var EnvMap = map[Environment][]string{
	EnvProduction:  []string{"production"},
	EnvDevelopment: []string{"development"},
	EnvTest:        []string{"test"},
}

// IsProduction returns true iff env is production.
func (e Environment) IsProduction() bool {
	return e == EnvProduction
}

// IsDevelopment returns true iff env is development.
func (e Environment) IsDevelopment() bool {
	return e == EnvDevelopment
}

// IsTest returns true iff env is test.
func (e Environment) IsTest() bool {
	return e == EnvTest
}

func (e Environment) String() string {
	s := EnvMap[e]
	if len(s) == 0 {
		return "unknown"
	}
	return s[0]
}

// GetEnvironment returns the app environment parsed from the
// environment variable.
func GetEnvironment() Environment {
	v := strings.ToLower(os.Getenv(EnvVarName))
	for env, l := range EnvMap {
		for _, e := range l {
			if e == v {
				return env
			}
		}
	}
	return EnvUnknown
}
