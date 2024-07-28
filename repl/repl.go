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

const PROMPT = "code>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprintf(out, PROMPT)
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
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
		if evaluated := evaluator.Eval(program, env); evaluated != nil {
			io.WriteString(out, evaluated.Inspect()+"\n")
		}
		/*
			tok := l.NextToken()
			for tok.Type != token.EOF {
				tok = l.NextToken()
				fmt.Fprintf(out, "%s\n", tok)
			}
		*/
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
