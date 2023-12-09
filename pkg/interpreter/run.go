package interpreter

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/taehioum/glox/pkg/parser"
	"github.com/taehioum/glox/pkg/scanner"
)

type Runner struct {
	// HadError bool
}

func (i *Runner) Runfile(path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("running file: %w", err)
	}

	return i.run(string(contents))
}

func (i *Runner) RunPrompt() error {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		b := sc.Scan()
		if !b {
			break
		}
		err := i.run(sc.Text())
		if err != nil {
			// return fmt.Errorf("running prompt: %w", err)
			fmt.Printf("running prompt %s: %s\n", sc.Text(), err)
		}
	}

	if sc.Err() != nil {
		return fmt.Errorf("running prompt: %w", sc.Err())
	}
	return nil
}

// the main logic
func (i *Runner) run(source string) error {
	tokens, err := scanner.ScanTokens(source)
	if err != nil {
		return fmt.Errorf("running: %w", err)
	}

	slog.Debug("tokens", slog.Attr{Key: "tokens", Value: slog.AnyValue(tokens)})
	stmts, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Errorf("running: %w", err)
	}

	slog.Debug("stmts", slog.Attr{Key: "stmts", Value: slog.AnyValue(stmts)})
	intpr := NewInterpreter(os.Stdout)
	err = intpr.Interprete(stmts...)
	if err != nil {
		return err
	}

	return nil
}
