package service

import (
	"os"
	"os/signal"
)

// Start is a convenience method equivalent to `service.Load(m).Run()` and starting the
// app with `./<myapp> start`. Prefer using `Run()` as it is more flexible.
func (s *Service) Start() error {
	return s.RunCommand("start")
}

// StartForTest starts the app with the environment set to test.
// Returns stop function as a convenience.
func (s *Service) StartForTest() func() {
	s.Env = EnvTest
	err := s.setup()
	if err != nil {
		panic(err)
	}
	go s.start()
	return s.Stop
}

// start calls Start on each module, in goroutines. Assumes that
// setup() has already been called. Blocks.
func (s *Service) start() {
	for _, m := range s.modules {
		n := getModuleName(m)
		c := s.configs[n]
		BootPrintln("[service] starting", n)
		if c.Start != nil {
			go c.Start()
		}
	}
	s.wait()
}

// wait blocks until a signal is received, or the stopper channel is closed
func (s *Service) wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case sig := <-c:
		BootPrintln("[service] got signal:", sig)
	case <-s.stopper:
		BootPrintln("[service] app stop")
	}
	s.stop()
}
