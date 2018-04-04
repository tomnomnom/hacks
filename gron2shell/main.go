package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s := statementFromString(scanner.Text())
		f := formatStatement(s)
		if f == "" {
			continue
		}
		fmt.Println(f)
	}
}

func formatStatement(s statement) string {
	out := &bytes.Buffer{}

	// strip off the leading 'json' bare key
	if s[0].typ == typBare && s[0].text == "json" {
		s = s[1:]
	}

	// strip off the leading dots
	if s[0].typ == typDot || s[0].typ == typLBrace {
		s = s[1:]
	}

	for _, t := range s {
		switch t.typ {
		case typBare:
			out.WriteString(t.text)

		case typNumericKey:
			out.WriteString(t.text)

		case typQuotedKey:
			out.WriteString(t.text[1 : len(t.text)-1])

		case typDot:
			out.WriteString(t.text)

		case typLBrace:
			out.WriteRune('.')

		case typRBrace:
			// nothing

		case typEquals:
			out.WriteString(t.text)

		case typSemi:
			// nothing

		case typString:
			out.WriteString(t.text[1 : len(t.text)-1])

		case typNumber:
			out.WriteString(t.text)

		case typTrue:
			out.WriteString(t.text)

		case typFalse:
			out.WriteString(t.text)

		case typNull:
			out.WriteString(t.text)

		case typEmptyArray:
			// ignore line
			return ""

		case typEmptyObject:
			// ignore line
			return ""

		default:
			// Nothing
		}
	}

	return out.String()
}
