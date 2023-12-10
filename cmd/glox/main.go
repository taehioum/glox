package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/taehioum/glox/pkg/runner"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}

	slogLeveler := slog.LevelInfo
	logLevel := os.Getenv("LOG")
	if logLevel == "DEBUG" || logLevel == "debug" {
		slogLeveler = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLeveler,
	})))

	i := runner.Runner{}

	var err error
	if len(args) == 1 {
		path := args[0]
		err = i.Runfile(path)
	} else {
		err = i.RunPrompt()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(65)
	}
}
