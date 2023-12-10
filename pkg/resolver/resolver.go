package resolver

import (
	"fmt"

	"github.com/taehioum/glox/pkg/ast"
	"github.com/taehioum/glox/pkg/interpreter"
)

type Resolver struct {
	interpreter *interpreter.Interpreter
	envs        []map[string]bool
}

func New(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		envs:        make([]map[string]bool, 0),
	}
}

func (r *Resolver) Resolve(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		if err := r.ResolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) ResolveExpr(expr ast.Expr) (any, error) {
	return expr.Accept(r)
}

func (r *Resolver) ResolveStmt(stmt ast.Stmt) error {
	return stmt.Accept(r)
}

func (r *Resolver) BeginScope() {
	r.envs = append(r.envs, make(map[string]bool))
}

func (r *Resolver) ExitScope() {
	r.envs = r.envs[:len(r.envs)-1]
}

func (r *Resolver) Declare(name string) {
	if len(r.envs) == 0 {
		return
	}
	r.envs[len(r.envs)-1][name] = false
}

func (r *Resolver) Define(name string) {
	if len(r.envs) == 0 {
		return
	}
	r.envs[len(r.envs)-1][name] = true
}

// VisitBlock implements ast.StatementVistior.
func (r *Resolver) VisitBlock(b ast.Block) error {
	r.BeginScope()
	defer r.ExitScope()

	return r.Resolve(b.Stmts)
}

// VisitAssignment implements ast.ExpressionVisitor.
func (r *Resolver) VisitAssignment(a ast.Assignment) (any, error) {
	if _, err := r.ResolveExpr(a.Value); err != nil {
		return nil, err
	}
	if err := r.resolveLocal(a, a.Name.Lexeme); err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitBinary implements ast.ExpressionVisitor.
func (r *Resolver) VisitBinary(b ast.Binary) (any, error) {
	if _, err := r.ResolveExpr(b.Left); err != nil {
		return nil, err
	}
	if _, err := r.ResolveExpr(b.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitCall implements ast.ExpressionVisitor.
func (r *Resolver) VisitCall(c ast.Call) (any, error) {
	if _, err := r.ResolveExpr(c.Callee); err != nil {
		return nil, err
	}
	for _, arg := range c.Args {
		if _, err := r.ResolveExpr(arg); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// VisitGrouping implements ast.ExpressionVisitor.
func (r *Resolver) VisitGrouping(g ast.Grouping) (any, error) {
	return r.ResolveExpr(g.Expr)
}

// VisitLambda implements ast.ExpressionVisitor.
func (r *Resolver) VisitLambda(l ast.Lambda) (any, error) {
	r.BeginScope()
	defer r.ExitScope()
	for _, param := range l.Params {
		r.Declare(param.Lexeme)
		r.Define(param.Lexeme)
	}
	return nil, r.Resolve(l.Body)
}

// VisitLiteral implements ast.ExpressionVisitor.
func (r *Resolver) VisitLiteral(l ast.Literal) (any, error) {
	return nil, nil
}

// VisitLogical implements ast.ExpressionVisitor.
func (r *Resolver) VisitLogical(l ast.Logical) (any, error) {
	if _, err := r.ResolveExpr(l.Left); err != nil {
		return nil, err
	}
	if _, err := r.ResolveExpr(l.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitPostUnary implements ast.ExpressionVisitor.
func (r *Resolver) VisitPostUnary(p ast.PostUnary) (any, error) {
	if _, err := r.ResolveExpr(p.Left); err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitUnary implements ast.ExpressionVisitor.
func (r *Resolver) VisitUnary(u ast.Unary) (any, error) {
	if _, err := r.ResolveExpr(u.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitVariable implements ast.ExpressionVisitor.
func (r *Resolver) VisitVariable(v ast.Variable) (any, error) {
	if len(r.envs) > 0 {
		b, ok := r.envs[len(r.envs)-1][v.Name.Lexeme]
		if ok && !b {
			return nil, fmt.Errorf("cannot read local variable in its own initializer")
		}
	}
	r.resolveLocal(v, v.Name.Lexeme)
	return nil, nil
}

func (r *Resolver) resolveLocal(expr ast.Expr, name string) error {
	for i := len(r.envs) - 1; i >= 0; i-- {
		if _, ok := r.envs[i][name]; ok {
			r.interpreter.Resolve(expr, len(r.envs)-1-i)
			return nil
		}
	}
	return nil
}

// VisitBreak implements ast.StatementVistior.
func (*Resolver) VisitBreak(ast.Break) error {
	return nil
}

// VisitContinue implements ast.StatementVistior.
func (*Resolver) VisitContinue(ast.Continue) error {
	return nil
}

// VisitDeclaration implements ast.StatementVistior.
func (r *Resolver) VisitDeclaration(decl ast.Declaration) error {
	r.Declare(decl.Name.Lexeme)
	if decl.Initializer != nil {
		if _, ok := decl.Initializer.(ast.Lambda); ok {
			// to allow recursive lambdas, we need to declare the variable first
			r.Define(decl.Name.Lexeme)
			_, err := r.ResolveExpr(decl.Initializer)
			return err
		}
		if _, err := r.ResolveExpr(decl.Initializer); err != nil {
			return err
		}
	}

	r.Define(decl.Name.Lexeme)
	return nil
}

// VisitExpression implements ast.StatementVistior.
func (r *Resolver) VisitExpression(e ast.Expression) error {
	_, err := r.ResolveExpr(e.Expr)
	return err
}

// VisitIf implements ast.StatementVistior.
func (r *Resolver) VisitIf(i ast.If) error {
	if _, err := r.ResolveExpr(i.Cond); err != nil {
		return err
	}
	if err := r.ResolveStmt(i.Then); err != nil {
		return err
	}
	if i.Else == nil {
		return nil
	}

	if err := r.ResolveStmt(i.Else); err != nil {
		return err
	}
	return nil
}

// VisitReturn implements ast.StatementVistior.
func (r *Resolver) VisitReturn(ret ast.Return) error {
	if ret.Value == nil {
		return nil
	}

	if _, err := r.ResolveExpr(ret.Value); err != nil {
		return err
	}
	return nil
}

// VisitWhile implements ast.StatementVistior.
func (r *Resolver) VisitWhile(w ast.While) error {
	if _, err := r.ResolveExpr(w.Cond); err != nil {
		return err
	}
	if err := r.ResolveStmt(w.Body); err != nil {
		return err
	}
	return nil
}

var _ ast.ExpressionVisitor = (*Resolver)(nil)
var _ ast.StatementVistior = (*Resolver)(nil)
