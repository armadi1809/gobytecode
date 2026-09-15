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
