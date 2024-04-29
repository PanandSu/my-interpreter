package parser

import (
	"fmt"
	"my-interpreter/ast"
	"my-interpreter/lexer"
	"my-interpreter/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var precedences = map[tkt]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESS_GREATER,
	token.GT:       LESS_GREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

const (
	_ int = iota
	LOWEST
	EQUALS
	LESS_GREATER
	SUM
	PRODUCT
	PREFIX
	INDEX
	CALL
)

type tkt = token.TokenType

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[tkt]prefixParseFn
	infixParseFns  map[tkt]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		errors:         []string{},
		prefixParseFns: map[tkt]prefixParseFn{},
		infixParseFns:  map[tkt]infixParseFn{},
	}
	{
		p.registerPrefix(token.IDENT, p.parseIdentifier)
		p.registerPrefix(token.INT, p.parseIntegerIdentifier)

		//注册布尔字面量
		p.registerPrefix(token.TRUE, p.parseBoolLiteral)
		p.registerPrefix(token.FALSE, p.parseBoolLiteral)

		//注册字符串字面量
		p.registerPrefix(token.STRING, p.parseStrLiteral)

		//注册数组字面量
		p.registerPrefix(token.LBRACKET, p.parseArrLiteral)

		//注册映射字面量
		p.registerPrefix(token.LBRACE, p.parseMapLiteral)

		//加了小括号的表达式
		p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

		//if表达式
		p.registerPrefix(token.IF, p.parseIfExpression)

		//函数定义字面量
		p.registerPrefix(token.FUNCTION, p.parseFuncLiteral)

		//前缀表达式
		p.registerPrefix(token.BANG, p.parsePrefixExpression)
		p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	}

	{
		//注册中缀表达式
		p.registerInfix(token.EQ, p.parseInfixExpression)
		p.registerInfix(token.NEQ, p.parseInfixExpression)
		p.registerInfix(token.LT, p.parseInfixExpression)
		p.registerInfix(token.GT, p.parseInfixExpression)
		p.registerInfix(token.PLUS, p.parseInfixExpression)
		p.registerInfix(token.MINUS, p.parseInfixExpression)
		p.registerInfix(token.ASTERISK, p.parseInfixExpression)
		p.registerInfix(token.SLASH, p.parseInfixExpression)

		//注册调用表达式
		p.registerInfix(token.LPAREN, p.parseCallExpression)

		//注册索引表达式
		p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	}

	//p.registerInfix(token.LPAREN, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t tkt) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tkt) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t tkt) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(typ tkt, fn prefixParseFn) {
	p.prefixParseFns[typ] = fn
}

func (p *Parser) registerInfix(typ tkt, fn infixParseFn) {
	p.infixParseFns[typ] = fn
}

// 开始解析程序
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	typ := p.curToken.Type
	switch typ {
	case token.SEMICOLON:
		return nil
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := ast.ExpressionStatement{}
	stmt.Expr = p.parseExpression(LOWEST)
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return &stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken}
}

func (p *Parser) parseIntegerIdentifier() ast.Expression {
	lit := ast.IntLiteral{Token: p.curToken}
	i, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("count not parse %q as an integer\n", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = i
	return &lit
}

func (p *Parser) parseExpressionListUntil(end token.TokenType) []ast.Expression {
	var exprs []ast.Expression
	if p.peekTokenIs(end) {
		p.nextToken()
		return exprs
	}
	p.nextToken()
	exprs = append(exprs, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exprs = append(exprs, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return exprs
}

// 解析字面量,包括布尔,字符串,数组,映射,函数
func (p *Parser) parseBoolLiteral() ast.Expression {
	return &ast.BoolLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseStrLiteral() ast.Expression {
	str := &ast.StrLiteral{}
	str.Token = p.curToken
	return str
}

func (p *Parser) parseArrLiteral() ast.Expression {
	arr := &ast.ArrLiteral{}
	p.parseExpressionListUntil(token.RBRACKET)
	return arr
}

func (p *Parser) parseMapLiteral() ast.Expression {
	m := ast.MapLiteral{}
	p.nextToken()
	key := p.parseExpression(LOWEST)
	if !p.expectPeek(token.COLON) {
		return nil
	}
	p.nextToken()
	val := p.parseExpression(LOWEST)
	m.Pairs[key] = val
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return &m
	}
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		val := p.parseExpression(LOWEST)
		m.Pairs[key] = val
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return &m
}

func (p *Parser) parseFuncLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fn.Parameters = p.parseFuncParameters()
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}
	/*if !p.expectPeek(token.LBRACE) {
		return nil
	}*/
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) || !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFuncParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	//无参数
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	ident := &ast.Identifier{Token: p.curToken}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident = &ast.Identifier{Token: p.curToken}
		identifiers = append(identifiers, ident)
	}
	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}
	return identifiers
}

// 解析表达式
// 包括前缀,后缀,if,调用,索引,分组
func (p *Parser) parsePrefixExpression() ast.Expression {
	prefix := ast.PrefixExpression{Token: p.curToken}
	p.nextToken()
	prefix.Right = p.parseExpression(PREFIX)
	return &prefix
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := ast.InfixExpression{Token: p.curToken, Left: left}
	p.nextToken()
	precedence := p.curPrecedence()
	expr.Right = p.parseExpression(precedence)
	return &expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	/*if !p.expectPeek(token.LBRACE) {
		return nil
	}
	p.nextToken()*/
	expr.Consequence = p.parseBlockStatement()
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		//if !p.expectPeek(token.LBRACE) {
		//	return nil
		//}
		//p.nextToken()
		expr.Alternative = p.parseBlockStatement()
		//if !p.expectPeek(token.RBRACE) {
		//	return nil
		//}
	}
	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.CallExpression{}
	call.Function = function
	call.Arguments = p.parseExpressionListUntil(token.RPAREN)
	return call
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	var index *ast.IndexExpression
	index.Left = left
	p.nextToken()
	index.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return index
}

/*func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	return args
}
*/
