package interpreter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/taehioum/glox/pkg/scanner"
)

type Interpreter struct {
	// HadError bool
}

func (i *Interpreter) Runfile(path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("running file: %w", err)
	}

	return i.run(string(contents))
}

func (i *Interpreter) RunPrompt() error {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		b := sc.Scan()
		if !b {
			break
		}
		fmt.Println(sc.Text())
		err := i.run(sc.Text())
		if err != nil {
			return fmt.Errorf("running prompt: %w", err)
		}
	}

	if sc.Err() != nil {
		return fmt.Errorf("running prompt: %w", sc.Err())
	}
	return nil
}

// the main logic
func (i *Interpreter) run(source string) error {
	tokens, err := scanner.ScanTokens(source)
	if err != nil {
		return fmt.Errorf("running: %w", err)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return nil
}
