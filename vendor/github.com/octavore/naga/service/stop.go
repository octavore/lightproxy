package service

import (
	"time"
)

// Stop the app by closing the stop channel. Used for testing.
func (s *Service) Stop() {
	close(s.stopper)

	// wait for all running modules to stop. todo: timeout?
	s.running.Wait()
}

// stop in reverse topological order
func (s *Service) stop() {
	dc := make(chan struct{})
	go func() {
		for i := len(s.modules) - 1; i >= 0; i-- {
			n := getModuleName(s.modules[i])
			c := s.configs[n]
			if c.Stop != nil {
				BootPrintln("[service] stopping", n)
				c.Stop()
			}
		}
		dc <- struct{}{}
	}()

	select {
	case <-dc:
	case <-time.After(30 * time.Second): // todo: make this configurable
		BootPrintln("[service] stop timed out")
	}
}
