package expr

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"time"
	"unicode"
)

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (l *lexer) next() {
	l.token = l.scan.Scan()
}

func (l *lexer) peek() rune {
	return l.scan.Peek()
}

func (l *lexer) text() string {
	return l.scan.TokenText()
}

func (l *lexer) has() bool {
	return l.token != scanner.EOF
}

func Parse(source string) (Expr, error) {
	if source == "" {
		return unary{true}, nil
	}
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(source))
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanFloats | scanner.ScanInts

	lex.next()
	e, err := parse(lex)
	if err != nil {
		return nil, err
	}
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected end of input: %s (%s)", string(lex.token), lex.text())
	}
	return e, nil
}

func parse(lex *lexer) (Expr, error) {
	lhs, err := parseExpr(lex)
	if err != nil {
		return nil, err
	}
	lex.next()
	if op := parseOp(lex); op == all || op == any {
		if rhs, err := parse(lex); err != nil {
			return nil, err
		} else {
			l := &logical{x: lhs, y: rhs, op: op}
			lex.next()
			return l, nil
		}
	} else {
		return lhs, nil
	}
}

func parseExpr(lex *lexer) (Expr, error) {
	lhs := parseValue(lex)
	lex.next()

	switch op := parseOp(lex); op {
	case eq, ne, gt, ge, lt, le, al, sw, ew:
		rhs := parseValue(lex)
		if rhs == invalid {
			return nil, NotFoundErr
		}
		return &binary{x: lhs, y: rhs, op: op}, nil
	default:
		return nil, fmt.Errorf("unsupported comparison expression: %s (%s)", UnsupportedOpErr, op)
	}
}

func parseValue(lex *lexer) Value {
	switch lex.token {
	case scanner.Ident:
		dtype := Default

		v := lex.text()
		if t := lex.peek(); t == ':' {
			lex.next()
			if v, ok := dtypes[lex.text()]; ok {
				dtype = v
			} else {
				dtype = Varchar
			}
		}
		return Var{name: v, dtype: dtype}
	case scanner.String:
		v := strings.TrimFunc(lex.text(), func(r rune) bool {
			return r == '"'
		})
		if d, err := time.Parse(time.RFC3339, v); err == nil {
			return datetime(d)
		} else {
			return varchar(v)
		}
	case scanner.Int, scanner.Float:
		v, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			v = 0.0
		}
		return number(v)
	default:
		return invalid
	}
}

func parseOp(lex *lexer) string {
	defer lex.next()
	switch lex.token {
	case '=':
		lex.scan.Next()
		return eq
	case '!':
		lex.scan.Next()
		return ne
	case '>', '<':
		runes := []rune{lex.token}
		next := lex.scan.Peek()
		if next == '=' {
			lex.scan.Next()
		}
		if !unicode.IsSpace(next) {
			runes = append(runes, next)
		}
		return string(runes)
	case '&':
		lex.scan.Next()
		return all
	case '|':
		lex.scan.Next()
		return any
	case '^':
		lex.scan.Next()
		return sw
	case '$':
		lex.scan.Next()
		return ew
	case '~':
		lex.scan.Next()
		return al
	default:
		return "??"
	}
}
