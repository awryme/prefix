package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/awryme/prefix/app/prefix"
	"github.com/awryme/prefix/pkg/cmdscanner"
)

const QuitCmd = ":q"

const desc = `prefix command.`

type App struct {
	Debug bool `short:"d" help:"print full command before executing"`

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

		// pass empty args by default
		var args []string
		if text != "" {
			args = strings.Split(text, " ")
		}

		err := executor.Run(ctx, args)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return cmdscanner.Stop
			}
			fmt.Println("executing error:", err)
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
