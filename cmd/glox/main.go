package main

import (
	"fmt"
	"os"

	"github.com/taehioum/glox/pkg/interpreter"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}

	i := interpreter.Interpreter{}

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
