package repl

import (
	"bufio"
	"fmt"
	"io"
	evaluator "my-interpreter/evaluator" //不要动
	"my-interpreter/lexer"
	"my-interpreter/object"
	"my-interpreter/parser"
)

const PROMPT = "JueJin>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		_, _ = fmt.Fprintf(out, PROMPT)
		if ok := scanner.Scan(); !ok {
			return
		}
		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		_, _ = io.WriteString(out, program.String())
		_, _ = io.WriteString(out, "\n")
		if evaluated := evaluator.Eval(program, env); evaluated != nil {
			_, _ = io.WriteString(out, evaluated.Inspect()+"\n")
		}
		/*
			tok := l.NextToken()
			for tok.Type != token.EOF {
				tok = l.NextToken()
				_, _ = fmt.Fprintf(out, "%s\n", tok)
			}
		*/
	}
}

const MONKEY_FACE = `            
   .--.  .-"__,__"-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, MONKEY_FACE)
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
