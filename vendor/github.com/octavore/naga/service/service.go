package service

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		BootPrintln = log.Println
	}
}

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

	defaultCommand string

	started sync.WaitGroup
	running sync.WaitGroup
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
	svc := loadEnv(m, GetEnvironment())
	svc.commands["help"] = &Command{
		Keyword: "help <command>",
		Run: func(ctx *CommandContext) {
			ctx.RequireExactlyNArgs(1)
			cmd := svc.commands[ctx.Args[0]]
			if cmd == nil {
				fmt.Printf("Unknown command: %s\n", ctx.Args[0])
				return
			}
			fmt.Printf("Usage of %s %s\n", os.Args[0], cmd.Keyword)
			if cmd.Usage == "" {
				fmt.Println(cmd.ShortUsage)
			} else {
				fmt.Println(cmd.Usage)
			}
		},
		ShortUsage: "Additional info for <command>",
		Usage:      "Show additional info for <command>",
	}
	return svc
}

// Run is a convenience method equivalent to "New(...).Run()"
func Run(m Module) {
	New(m).Run()
}

// Load the app with the given environment, and initializes
// all modules recursively starting with m.
func loadEnv(m Module, env Environment) *Service {
	BootPrintln("[service] env is", env.String())
	svc := &Service{
		Env:      env,
		stopper:  make(chan bool),
		modules:  []Module{},
		configs:  map[string]*Config{},
		commands: map[string]*Command{},
		started:  sync.WaitGroup{},
	}
	svc.started.Add(1)
	svc.load(m)
	return svc
}

// Usage prints the usage for all registered commands.
func (s *Service) Usage() {
	fmt.Printf("Usage of %s\n", os.Args[0])
	if s.commands["help"] != nil {
		fmt.Printf("    %-16s %s\n", "help", s.commands["help"].ShortUsage)
	}
	if s.commands["start"] != nil {
		fmt.Printf("    %-16s %s\n", "start", s.commands["start"].ShortUsage)
	}
	for k, cmd := range s.commands {
		if k == "start" || k == "help" {
			continue
		}
		fmt.Printf("    %-16s %s\n", k, cmd.ShortUsage)
	}
}

// Run parses arguments from the command line and passes them to RunCommand.
func (s *Service) Run() {
	flag.Usage = s.Usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 && s.defaultCommand != "" {
		args = append([]string{s.defaultCommand}, args...)
	}
	if len(args) == 0 {
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
