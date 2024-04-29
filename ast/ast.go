package ast

import (
	"bytes"
	"my-interpreter/token"
	"strings"
)

type Node interface {
	String() string
}

// 只要实现了String()方法,就是Expression
type Expression interface {
	Node
}

// 只要实现了String()方法,就是Statement
type Statement interface {
	Node
}

// 程序
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	//bingbing!
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// 标识符
type Identifier struct {
	Token token.Token
}

func (i *Identifier) String() string {
	return i.Token.Literal
}

// 整数字面量
type IntLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntLiteral) String() string {
	return il.Token.Literal
}

// 布尔字面量
type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BoolLiteral) String() string {
	return bl.Token.Literal
}

// 字符串字面量
type StrLiteral struct {
	Token token.Token
}

func (sl *StrLiteral) String() string {
	return sl.Token.Literal
}

// 数组字面量
type ArrLiteral struct {
	Elements []Expression
}

func (al *ArrLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	var elements []string
	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// 映射字面量
type MapLiteral struct {
	Pairs map[Expression]Expression
}

func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	var pairs []string
	for key, value := range ml.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// 实现了String()方法,隐式实现了Statement接口
type LetStatement struct {
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString("let")
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}

// 返回语句
type ReturnStatement struct {
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("return")
	out.WriteString(" ")
	out.WriteString(rs.ReturnValue.String())
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}

type ExpressionStatement struct {
	Expr Expression
}

func (es *ExpressionStatement) String() string {
	return es.Expr.String() + ";" + "\n"
}

// 前缀表达式
type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Token.Literal)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// 中缀表达式
type InfixExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Token.Literal + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// if表达式
type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" ")
		out.WriteString("else")
		out.WriteString(" ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// 语句块
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	out.WriteString("\n")
	for _, s := range bs.Statements {
		out.WriteString("\t" + s.String())
	}
	out.WriteString("}")
	return out.String()
}

// 函数字面量
type FunctionLiteral struct {
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// 函数调用表达式
type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	var args []string
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// 索引表达式
type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}
