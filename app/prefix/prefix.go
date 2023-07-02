package prefix

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Executor struct {
	// print full command to passed stderr before running
	PrintOnRun bool

	binary   string
	initArgs []string

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewExecutor(stdin io.Reader, stdout, stderr io.Writer, binary string, args []string) *Executor {
	return &Executor{
		binary:   binary,
		initArgs: args,
		stdin:    stdin,
		stdout:   stdout,
		stderr:   stderr,
	}
}

func (e *Executor) Run(ctx context.Context, args []string) error {
	fullArgs := append(e.initArgs, args...)

	e.printArgs(fullArgs)

	cmd := exec.CommandContext(ctx, e.binary, fullArgs...)
	cmd.Stdin = e.stdin
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed command: %w", err)
	}
	return nil
}

func (e *Executor) printArgs(fullArgs []string) {
	if e.PrintOnRun {
		fmt.Fprintln(e.stderr, "running:", e.binary, strings.Join(fullArgs, " "))
	}
}
