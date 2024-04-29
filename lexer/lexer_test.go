package lexer

import (
	"my-interpreter/token"
	"testing"
)

type Expect struct {
	Type    token.TokenType
	Literal string
}

func NewExpect(typ token.TokenType, literal string) Expect {
	return Expect{typ, literal}
}

const testData = `
let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`

var tests = []Expect{
	let, five, assign, int_5, semicolon,
	let, ten, assign, int_10, semicolon,
	let, add, assign, fn, lparen, x, comma, y, rparen, lbrace,
	x, plus, y, semicolon,
	rbrace, semicolon,
	let, result, assign, add, lparen, five, comma, ten, rparen, semicolon,
	bang, minus, slash, asterisk, int_5, semicolon,
	int_5, lt, int_10, gt, int_5, semicolon,
	if_, lparen, int_5, lt, int_10, rparen, lbrace,
	return_, true_, semicolon,
	rbrace, else_, lbrace,
	return_, false_, semicolon,
	rbrace,
	int_10, eq, int_10, semicolon,
	int_10, neq, int_9, semicolon,
	str_foobar,
	str_foo_bar,
	lbracket, int_1, comma, int_2, rbracket, semicolon,
	lbrace, str_foo, colon, str_bar, rbrace,
	eof,
}

func TestLexer_NextToken(t *testing.T) {
	input := testData
	l := NewLexer(input)
	for i, expect := range tests {
		tok := l.NextToken()
		if tok.Type != expect.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, expect.Type, tok.Type)
		}
		if tok.Literal != expect.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, expect.Literal, tok.Literal)
		}
	}
}

var (
	// testData中用到的identifier
	five   = NewExpect(token.IDENT, "five")
	ten    = NewExpect(token.IDENT, "ten")
	add    = NewExpect(token.IDENT, "add")
	x      = NewExpect(token.IDENT, "x")
	y      = NewExpect(token.IDENT, "y")
	result = NewExpect(token.IDENT, "result")

	eof    = NewExpect(token.EOF, "")
	int_5  = NewExpect(token.INT, "5")
	int_10 = NewExpect(token.INT, "10")
	int_9  = NewExpect(token.INT, "9")
	int_1  = NewExpect(token.INT, "1")
	int_2  = NewExpect(token.INT, "2")

	str_foobar  = NewExpect(token.STRING, "foobar")
	str_foo_bar = NewExpect(token.STRING, "foo bar")
	str_foo     = NewExpect(token.STRING, "foo")
	str_bar     = NewExpect(token.STRING, "bar")

	assign   = NewExpect(token.ASSIGN, "=")
	plus     = NewExpect(token.PLUS, "+")
	minus    = NewExpect(token.MINUS, "-")
	asterisk = NewExpect(token.ASTERISK, "*")
	slash    = NewExpect(token.SLASH, "/")

	bang = NewExpect(token.BANG, "!")
	lt   = NewExpect(token.LT, "<")
	gt   = NewExpect(token.GT, ">")

	eq  = NewExpect(token.EQ, "==")
	neq = NewExpect(token.NEQ, "!=")

	comma     = NewExpect(token.COMMA, ",")
	colon     = NewExpect(token.COLON, ":")
	semicolon = NewExpect(token.SEMICOLON, ";")

	lparen   = NewExpect(token.LPAREN, "(")
	rparen   = NewExpect(token.RPAREN, ")")
	lbracket = NewExpect(token.LBRACKET, "[")
	rbracket = NewExpect(token.RBRACKET, "]")
	lbrace   = NewExpect(token.LBRACE, "{")
	rbrace   = NewExpect(token.RBRACE, "}")

	//testDate中用到的关键字
	let     = NewExpect(token.LET, "let")
	fn      = NewExpect(token.FUNCTION, "fn")
	true_   = NewExpect(token.TRUE, "true")
	false_  = NewExpect(token.FALSE, "false")
	if_     = NewExpect(token.IF, "if")
	else_   = NewExpect(token.ELSE, "else")
	return_ = NewExpect(token.RETURN, "return")
)
