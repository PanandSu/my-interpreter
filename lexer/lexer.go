package lexer

import "my-interpreter/token"

type Lexer struct {
	input     string
	index     int
	nextIndex int
	char      byte
}

func NewLexer(input string) *Lexer {
	l := Lexer{input: input, index: -1, nextIndex: 0}
	l.readChar()
	return &l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.char {
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '=':
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			s := string(ch) + string(l.char)
			tok.Type = token.EQ
			tok.Literal = s
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case '<':
		tok = newToken(token.LT, l.char)
	case '!':
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			s := string(ch) + string(l.char)
			tok.Type = token.NEQ
			tok.Literal = s
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ':':
		tok = newToken(token.COLON, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '[':
		tok = newToken(token.LBRACKET, l.char)
	case ']':
		tok = newToken(token.RBRACKET, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if token.IsLetter(l.char) {
			s := l.readIdentifier()
			tok.Literal = s
			//判断是关键字还是标识符
			tok.Type = token.LookupIdentifier(s)
			return tok
		} else if token.IsDigit(l.char) {
			s := l.readNumber()
			tok.Literal = s
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() byte {
	l.index++
	if l.nextIndex >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.index]
	}
	l.nextIndex++
	return l.char
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func (l *Lexer) readIdentifier() string {
	index := l.index
	for token.IsLetter(l.char) {
		l.readChar()
	}
	return l.input[index:l.index]
}

func (l *Lexer) skipWhitespace() {
	//注意是for循环!
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	index := l.index
	for token.IsDigit(l.char) {
		l.readChar()
	}
	return l.input[index:l.index]
}

func (l *Lexer) readString() string {
	l.readChar()
	begin := l.index
	for l.char != '"' && l.char != 0 {
		l.readChar()
	}
	//跳出循环时l.char已经是字符串下引号了
	return l.input[begin:l.index]
}

func (l *Lexer) peekChar() byte {
	if l.index >= len(l.input) {
		return 0
	} else {
		//只是窥视,并未读取
		return l.input[l.nextIndex]
	}
}
