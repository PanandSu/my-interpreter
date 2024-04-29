package ast

import (
	"fmt"
	"my-interpreter/token"
	"testing"
)

func TestLetStatement_String(t *testing.T) {
	program := &Program{Statements: []Statement{
		&LetStatement{
			Name: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "my-interpreter",
				},
			},
			Value: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "panjinhao",
				},
			},
		},
	}}
	if program.String() != "let my-interpreter = panjinhao;\n" {
		t.Errorf("program.StrLiteral() wrong. got=%q", program.String())
	}
	fmt.Printf("%#v", program.String())
}

// 注意,不要修改fn的格式,开头没有\n,结尾也没有,你一旦按了Enter就会多加,就错了
var fn = `fn(a, b) {
	return (a + b);
}`

func TestFunctionLiteral_String(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&FunctionLiteral{
				Parameters: []*Identifier{
					{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "a",
						},
					},
					{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "b",
						},
					},
				},
				Body: &BlockStatement{
					Statements: []Statement{
						&ReturnStatement{
							ReturnValue: &InfixExpression{
								Token: token.Token{
									Type:    token.PLUS,
									Literal: "+",
								},
								Left: &IntLiteral{
									Token: token.Token{
										Type:    token.IDENT,
										Literal: "a",
									},
								},
								Right: &Identifier{
									Token: token.Token{
										Type:    token.IDENT,
										Literal: "b",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if program.String() != fn {
		t.Errorf("program.StrLiteral() wrong. got=%q", program.String())
	}
	fmt.Printf("%#v\n", program.String())
	fmt.Println(fn)
}

var if_stmt = `if (!a) {
	return true;
} else {
	return 2;
}`

func TestIfExpression_String(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&IfExpression{
				Condition: &PrefixExpression{
					Token: token.Token{
						Type:    token.BANG,
						Literal: "!",
					},
					Right: &Identifier{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "a",
						},
					},
				},
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ReturnStatement{
							ReturnValue: &BoolLiteral{
								Token: token.Token{
									Type:    token.TRUE,
									Literal: "true",
								},
								Value: true,
							},
						},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ReturnStatement{
							ReturnValue: &IntLiteral{
								Token: token.Token{
									Type:    token.INT,
									Literal: "2",
								},
								Value: 2,
							},
						},
					},
				},
			},
		},
	}
	if program.String() != if_stmt {
		t.Errorf("program.StrLiteral() wrong. got=%q", program.String())
	}
	fmt.Printf("%#v\n", program.String())
	fmt.Println(program.String())
	fmt.Printf("%#v\n", if_stmt)
}
