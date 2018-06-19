package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/tomnomnom/rplex"
)

func main() {
	flag.Parse()

	fn := flag.Arg(0)
	if fn == "" {
		fmt.Println("usage: jsstrings <filename>")
		return
	}

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := string(b)

	lex := rplex.New(s)
	tokens := lex.Run(lexText)

	for _, t := range tokens {
		if t.Text() == "" {
			continue
		}
		fmt.Printf("%s\n", t.Text())
	}
}

type stringLiteral struct {
	rplex.TextToken
}

type comment struct {
	rplex.TextToken
}

func lexText(l *rplex.Lexer) rplex.LexFn {
	l.AcceptUntil("'\"`/")
	l.Ignore()
	r := l.Peek()

	switch r {
	case '"':
		return lexDouble
	case '\'':
		return lexSingle
	case '`':
		return lexBacktick
	case '/':
		return lexComment
	default:
		return nil
	}

	return nil
}

func lexDouble(l *rplex.Lexer) rplex.LexFn {
	l.Accept("\"")
	l.Ignore()

	l.AcceptUntilUnescaped("\"")
	l.Emit(&stringLiteral{})

	l.Accept("\"")
	l.Ignore()

	return lexText
}

func lexSingle(l *rplex.Lexer) rplex.LexFn {
	l.Accept("'")
	l.Ignore()

	l.AcceptUntilUnescaped("'")
	l.Emit(&stringLiteral{})

	l.Accept("'")
	l.Ignore()

	return lexText
}

func lexBacktick(l *rplex.Lexer) rplex.LexFn {
	l.Accept("`")
	l.Ignore()

	l.AcceptUntilUnescaped("`")
	l.Emit(&stringLiteral{})

	l.Accept("`")
	l.Ignore()

	return lexText
}

func lexComment(l *rplex.Lexer) rplex.LexFn {
	l.Accept("/")

	line := l.Accept("/")
	block := l.Accept("*")
	l.Ignore()

	if !line && !block {
		// false alarm; division or a regex literal
		return lexToEndOfStatement
	}

	if line {
		l.AcceptUntil("\n")
		l.Emit(&comment{})
		return lexText
	}

	// read until we hit '*/'
	for {
		if l.Cur == utf8.RuneError {
			return nil
		}
		l.AcceptUntil("*")
		l.Accept("*")
		if l.Peek() == '/' {
			l.Accept("/")
			l.Emit(&comment{})
			break
		}
	}

	return lexText
}

func lexToEndOfStatement(l *rplex.Lexer) rplex.LexFn {
	l.AcceptUntilUnescaped(";\n")
	l.Accept(";\n")
	l.Ignore()
	return lexText
}
