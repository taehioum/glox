package interpreter

import (
	"errors"

	statements "github.com/taehioum/glox/pkg/ast"
	"github.com/taehioum/glox/pkg/interpreter/environment"
)

func (i *Interpreter) VisitDeclaration(stmt statements.Declaration) error {
	if stmt.Intializer == nil {
		i.env.Define(stmt.Name.Lexeme, nil)
		return nil
	}

	v, err := i.Eval(stmt.Intializer)
	if err != nil {
		return err
	}
	i.env.Define(stmt.Name.Lexeme, v)
	return nil
}

func (i *Interpreter) VisitBlock(stmt statements.Block) error {
	prev := i.env
	defer func() {
		// restore env
		i.env = prev
	}()
	i.env = environment.NewEnclosedEnvironment(prev)
	for _, stmt := range stmt.Stmts {
		err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitIf(stmt statements.If) error {
	v, err := i.Eval(stmt.Cond)
	if err != nil {
		return err
	}
	if truthy(v) {
		return stmt.Then.Accept(i)
	}
	if stmt.Else != nil {
		return stmt.Else.Accept(i)
	}
	return nil
}

func (i *Interpreter) VisitWhile(stmt statements.While) error {
	for {
		v, err := i.Eval(stmt.Cond)
		if err != nil {
			return err
		}
		if !truthy(v) {
			break
		}
		err = stmt.Body.Accept(i)
		if errors.Is(err, ErrBreak) {
			return nil
		}
		// by continuing the for loop here, we get the continue behavior for free.
		if errors.Is(err, ErrContinue) {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitBreak(stmt statements.Break) error {
	return ErrBreak
}

func (i *Interpreter) VisitContinue(stmt statements.Continue) error {
	return ErrContinue
}

func (i *Interpreter) VisitExpression(stmt statements.Expression) error {
	_, err := i.Eval(stmt.Expr)
	return err
}

func (i *Interpreter) VisitReturn(stmt statements.Return) error {
	if stmt.Value == nil {
		return ErrReturn{}
	}
	v, err := i.Eval(stmt.Value)
	if err != nil {
		return err
	}
	return ErrReturn{Value: v}
}
