package expression

type Expression interface {
	expr()
}

type IntExpr int

func (IntExpr) expr() {}

type ValExpr struct {
	Name string
	Expr Expression
}

func (ValExpr) expr() {}

type NameExpr string

func (NameExpr) expr() {}

type CallExpr struct {
	Function Expression
	Args     []Expression
}

func (CallExpr) expr() {}

type IfExpr struct {
	Cond    Expression
	IfTrue  Expression
	IfFalse Expression
}

func (IfExpr) expr() {}

type LambdaExpr struct {
	Params []string
	Body   Expression
}

func (LambdaExpr) expr() {}

type BeginExpr struct {
	Exps []Expression
}

func (BeginExpr) expr() {}

func Call(funcName string, args ...Expression) CallExpr {
	return CallExpr{
		Function: NameExpr(funcName),
		Args:     args,
	}
}

func Define(funcName string, params []string, body Expression) ValExpr {
	lambda := LambdaExpr{
		Params: params,
		Body:   body,
	}

	return ValExpr{
		Name: funcName,
		Expr: lambda,
	}
}

func If(cond Expression, ifT Expression, ifF Expression) IfExpr {
	return IfExpr{
		Cond:    cond,
		IfTrue:  ifT,
		IfFalse: ifF,
	}
}

func Begin(exps ...Expression) BeginExpr {
	return BeginExpr{Exps: exps}
}
