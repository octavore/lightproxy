package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Flag struct {
	Key   string
	Value *string
	Usage string
}

// Command represents a command-line keyword for the app.
// This is then typically invoked as follows:
//   ./myapp <keyword>
type Command struct {
	Keyword    string
	Run        func(*CommandContext)
	ShortUsage string
	Usage      string
	Flags      []*Flag
}

// AddCommand adds a command to the service via its Config.
func (c *Config) AddCommand(cmd *Command) {
	c.service.registerCommand(cmd)
}

// CommandContext is passed to the command when it is run,
// containing an array of parsed arguments.
type CommandContext struct {
	cmd   *Command
	Args  []string
	Flags map[string]*Flag
}

// UsageExit prints the usage for the executed command and exits.
func (c *CommandContext) UsageExit() {
	fmt.Println(c.cmd.Keyword)
	fmt.Println(c.cmd.Usage)
	os.Exit(0)
}

func (c *CommandContext) Fatal(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}

// RequireAtLeastNArgs is a helper function to ensure we have at least n args.
func (c *CommandContext) RequireAtLeastNArgs(n int) {
	if len(c.Args) < n {
		c.UsageExit()
	}
}

// RequireExactlyNArgs is a helper function to ensure we have at exactly n args.
func (c *CommandContext) RequireExactlyNArgs(n int) {
	if len(c.Args) != n {
		c.UsageExit()
	}
}

func parseArgs(flags []*Flag, args []string) (map[string]*Flag, []string, error) {
	flagMap := map[string]*Flag{}
	for _, f := range flags {
		ks := strings.Split(f.Key, ",")
		for _, k := range ks {
			k = strings.TrimSpace(k)
			k = strings.TrimLeft(k, "-")
			flagMap[k] = f
		}
	}

	updatedArgs := []string{}
	for i := 0; i < len(args); i++ {
		k := args[i]
		if !isFlag(k) {
			updatedArgs = append(updatedArgs, k)
			continue
		}
		k = strings.TrimLeft(k, "-")
		f, ok := flagMap[k]
		if !ok {
			return nil, nil, errors.New("no flag found: " + k)
		}
		if i+1 == len(args) || isFlag(args[i+1]) {
			f.Value = ptr("")
		} else {
			f.Value = ptr(args[i+1])
			i++
		}
	}
	return flagMap, updatedArgs, nil
}

func isFlag(f string) bool {
	return strings.HasPrefix(f, "-")
}

func ptr(s string) *string {
	return &s
}
