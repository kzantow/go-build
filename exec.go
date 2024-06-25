package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Run a command, logging with current stdout / stderr
func Run(cmd ...string) {
	cmd = append(Shelly(cmd[0]), cmd[1:]...)
	for i := range cmd {
		cmd[i] = Tpl(cmd[i])
	}
	NoErr(Exec(Path(cmd[0]), ExecArgs(cmd[1:]...), ExecStd()))
}

// Exec executes the given command, returning stdout and any error information
//
//nolint:gosec
func Exec(cmd Path, opts ...ExecOpt) error {
	// add the ToolDir first in the path for easier script writing
	lookupPath := os.Getenv("PATH")
	defer func() { LogErr(os.Setenv("PATH", lookupPath)) }()
	NoErr(os.Setenv("PATH", Tpl(ToolDir)+string(os.PathListSeparator)+lookupPath))

	// create the command, this will look it up based on path:
	c := exec.CommandContext(ctx, string(cmd))
	c.Env = os.Environ()
	for k, v := range Globals {
		val := ""
		switch v := v.(type) {
		case func() string:
			val = v()
		case string:
			val = Tpl(v)
		default:
			continue
		}
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", k, val))
	}

	// finally, apply all the options to modify the command
	for _, opt := range opts {
		opt(c)
	}

	args := c.Args[1:] // exec.Command sets the cmd to Args[0]
	Log("%v %v", cmd, strings.Join(args, " "))

	// execute
	err := c.Start()
	if err == nil {
		err = c.Wait()
	}
	if err != nil || (c.ProcessState != nil && c.ProcessState.ExitCode() > 0) {
		return &StackTraceError{
			Err:      fmt.Errorf("error executing: %s %s: %w", cmd, printArgs(args), err),
			ExitCode: c.ProcessState.ExitCode(),
		}
	}
	return nil
}

// ExecArgs appends args to the command
func ExecArgs(args ...string) ExecOpt {
	return func(cmd *exec.Cmd) {
		cmd.Args = append(cmd.Args, args...)
	}
}

// ExecStd executes with output mapped to the current process' stdout, stderr, stdin
func ExecStd() ExecOpt {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}
}

// ExecOut sends stdout to the writer
func ExecOut(out io.Writer) ExecOpt {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}
}

// ExecEnv adds an environment variable to the command
func ExecEnv(key, val string) ExecOpt {
	return func(cmd *exec.Cmd) {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, Tpl(val)))
	}
}

// ExecOpts combines all opts into a single one
func ExecOpts(opts ...ExecOpt) ExecOpt {
	return func(cmd *exec.Cmd) {
		for _, opt := range opts {
			opt(cmd)
		}
	}
}

// ExecOpt is used to alter the command used in Exec calls
type ExecOpt func(*exec.Cmd)

var ctx, cancel = context.WithCancel(context.Background())

// Cancel invokes the cancel call on all active
func Cancel() {
	cancel()
}

func printArgs(args []string) string {
	for i, arg := range args {
		if strings.Contains(arg, " ") {
			if strings.Contains(arg, `'`) {
				args[i] = `"` + arg + `"`
			} else {
				args[i] = "'" + arg + "'"
			}
		}
	}
	return strings.Join(args, " ")
}
