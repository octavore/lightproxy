package service

// Config contains functions for handling a Module's lifecycle.
// Each Module is responsible for setting up its Config object
// when its Init method is called (which can happen in any order).
// When a command is run, Setup is always called.
// If the command is the start command, Config.Start is also invoked.
type Config struct {
	// Setup is called regardless of the command, after initialization
	// and before the command (e.g., before Start).
	// Setup is executed in topological order sequentially (leaves first).
	Setup func() error

	// SetupTest is called after Setup and before Start (or command).
	// Always called in tests.
	SetupTest func()

	// Start is invoked when the app is run with '<myapp> start', i.e.
	// the start command. Start is run in goroutines, and is launched
	// sequentially in topological order.
	Start func()

	// Stop is invoked when a ctrl-c signal is received. May not execute,
	// e.g. if the app is kill-9ed. A maximum of 30 seconds is given
	// for all Stop functions to finish. Stop is invoked sequentially
	// in reverse topological order (parents first).
	Stop func()

	dependencies []Module // dependencies of this Module
	parent       Module   // parent of this Module
	service      *Service // pointer to the service
}

// Env returns the service environment.
func (c *Config) Env() Environment {
	return c.service.Env
}
