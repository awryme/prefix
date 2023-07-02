package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/anmitsu/go-shlex"
	"github.com/awryme/prefix/app/prefix"
	"github.com/awryme/prefix/pkg/cmdscanner"
)

const QuitCmd = ":q"

const desc = `prefix: execute commands in repl. Use ctrl+C or ":q" to quit`

type App struct {
	Debug bool `short:"d" help:"print full command to stderr before executing"`

	Binary string   `arg:"" help:"binary to run" passthrough:""`
	Args   []string `arg:"" help:"initial arguments" optional:""`
}

func Run() error {
	var app App
	cmdCtx := kong.Parse(&app, kong.UsageOnError(), kong.Name("prefix"), kong.Description(desc))
	cmdCtx.FatalIfErrorf(cmdCtx.Validate())

	ctx, cancel := context.WithCancel(context.Background())
	go cancelOnSignal(cancel)
	go exitOnDone(ctx)

	executor := prefix.NewExecutor(os.Stdin, os.Stdout, os.Stderr, app.Binary, app.Args)
	executor.PrintOnRun = app.Debug

	err := cmdscanner.Scan(os.Stdin, "> ", func(text string) error {
		if text == QuitCmd {
			return cmdscanner.Stop
		}

		args, err := shlex.Split(text, true)
		if err != nil {
			fmt.Println("parse args error:", err)
			return nil
		}

		err = executor.Run(ctx, args)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return cmdscanner.Stop
			}
			fmt.Println("execution error:", err)
		}
		return nil
	})

	return err
}

func cancelOnSignal(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// wait and cancel
	<-c
	cancel()
}

func exitOnDone(ctx context.Context) {
	<-ctx.Done()
	os.Exit(0)
}

func main() {
	err := Run()
	if err != nil {
		fmt.Println("exiting on error:", err)
		os.Exit(1)
	}
}
