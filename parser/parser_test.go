package parser

import (
	"my-interpreter/lexer"
	"testing"
)

func TestLiteralExpression(t *testing.T) {}

func TestIfExpression(t *testing.T) {}

func TestIfElseExpression(t *testing.T) {}

func TestReturnStatement(t *testing.T) {}

func TestIdentifierExpression(t *testing.T) {}

func TestLetStatement(t *testing.T) {}

func TestParsingPrefixExpression(t *testing.T) {}

func TestParsingInfixExpression(t *testing.T) {}

func TestCallExpressionParsing(t *testing.T) {}

func TestFunctionLiteralParsing(t *testing.T) {}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b);\n",
		},
		{
			"!-a",
			"(!(-a));\n",
		},
		{
			"a + b + c",
			"((a + b) + c);\n",
		},
		{
			"a + b - c",
			"((a + b) - c);\n",
		},
		{
			"a * b * c",
			"((a * b) * c);\n",
		},
		{
			"a * b / c",
			"((a * b) / c);\n",
		},
		{
			"a + b / c",
			"(a + (b / c));\n",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f);\n",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4);\n((-5) * 5);\n",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4));\n",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4));\n",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));\n",
		},
		{
			"true",
			"true;\n",
		},
		{
			"false",
			"false;\n",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false);\n",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true);\n",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4);\n",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2);\n",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5));\n",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5));\n",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5));\n",
		},
		{
			"!(true == true)",
			"(!(true == true));\n",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d);\n",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));\n",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g));\n",
		},
	}
	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
