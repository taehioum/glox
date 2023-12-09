package tests

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taehioum/glox/pkg/interpreter"
)

//go:embed fib.lox
var fib string

func TestFib(t *testing.T) {
	r := interpreter.Runner{}
	var b bytes.Buffer
	err := r.Run(fib, io.Writer(&b))
	if err != nil {
		t.Fatalf("running fib.lox: %s", err)
	}
	assert.Equal(t, "0\n1\n1\n2\n3\n5\n8\n13\n21\n34\n", b.String())
}
