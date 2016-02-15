package service

import (
	"flag"
	"fmt"
	"strings"
)

// BootPrintln can be replaced with log.Println for printing debug information.
var BootPrintln = func(v ...interface{}) {}

// A Service wraps a Module (and its dependencies). It manages the lifecycle
// and allows the Module to be started and stopped. It maintains a
// topologically sorted list of the Modules, along with a map of the Modules'
// Configs and a map of registered commands.
type Service struct {
	Env Environment

	stopper  chan bool
	modules  []Module
	configs  map[string]*Config
	commands map[string]*Command
}

// New creates a new service with Module m as the entry point
func New(m Module) *Service {
	svc := NewApp(m)
	svc.commands["start"] = &Command{
		Keyword:    "start",
		Run:        func(*CommandContext) { svc.start() },
		ShortUsage: "Start the app",
		Usage:      "Start running the app",
	}
	return svc
}

// NewApp creates a new app with Module m as the entry point. Unlike
// New, `start` is not automatically registered.
func NewApp(m Module) *Service {
	return loadEnv(m, GetEnvironment())
}

// Run is a convenience method equivalent to "New(...).Run()"
func Run(m Module) {
	New(m).Run()
}

// Load the app with the given environment, and initializes
// all modules recursively starting with m.
func loadEnv(m Module, env Environment) *Service {
	svc := &Service{
		Env:      env,
		stopper:  make(chan bool),
		modules:  []Module{},
		configs:  map[string]*Config{},
		commands: map[string]*Command{},
	}
	svc.load(m)
	return svc
}

// Usage prints the usage for all registered commands.
func (s *Service) Usage() {
	for k, cmd := range s.commands {
		fmt.Printf("    %-16s %s\n", k, cmd.ShortUsage)
	}
}

// Run parses arguments from the command line and passes them to RunCommand.
func (s *Service) Run() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		s.Usage()
		BootPrintln()
		return
	}
	err := s.RunCommand(args[0], args[1:]...)
	if err != nil {
		panic(err)
	}
}

// RunCommand executes the given command, or returns an error if not found.
// module setup (and setupTest) will be called recursively before
// executing the command via cmd.Run. Meant for tests.
func (s *Service) RunCommand(command string, args ...string) error {
	cmd := s.commands[command]
	if cmd == nil {
		return fmt.Errorf("unknown command %q", command)
	}
	err := s.setup()
	if err != nil {
		return fmt.Errorf("error in setup: %v", err)
	}
	flagMap, args, err := parseArgs(cmd.Flags, args)
	if err != nil {
		return err
	}
	cmd.Run(&CommandContext{cmd, args, flagMap})
	return nil
}

// setup invokes `Setup()` on all loaded modules in topological order,
// dependencies first. If `service.Env.IsTest()`, it also runs
// each module's `SetupTest()` immediately after the module's `Setup()`
func (s *Service) setup() error {
	for _, m := range s.modules {
		n := getModuleName(m)
		c := s.configs[n]
		if c.Setup != nil {
			BootPrintln("[service] setup", n)
			err := c.Setup()
			if err != nil {
				return err
			}
		}

		if s.Env.IsTest() && c.SetupTest != nil {
			BootPrintln("[service] setup for test", n)
			c.SetupTest()
		}
	}
	return nil
}

func (s *Service) getConfig(m Module) *Config {
	n := getModuleName(m)
	return s.configs[n]
}

func (s *Service) registerCommand(cmd *Command) {
	kw := strings.Split(cmd.Keyword, " ")[0]
	_, ok := s.commands[kw]
	if ok {
		panic("keyword already registered: " + kw)
	}
	s.commands[kw] = cmd
}
