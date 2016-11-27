package expr

import (
	"fmt"
)

type unary struct {
	value bool
}

func (u unary) Compare(e Env) bool {
	return u.value
}

func (u unary) String() string {
	return fmt.Sprintf("%t", u.value)
}

type binary struct {
	x, y Value
	op   string
}

func (b *binary) Compare(e Env) bool {
	x := b.x.Lookup(e)
	y := b.y.Lookup(e)
	switch b.op {
	case eq:
		return x.Equal(y)
	case ne:
		return !x.Equal(y)
	case lt:
		return !x.Equal(y) && x.Less(y)
	case le:
		return x.Equal(y) || x.Less(y)
	case gt:
		return !x.Equal(y) && !x.Less(y)
	case ge:
		return x.Equal(y) || !x.Less(y)
	case al:
		return x.Almost(y)
	case sw:
		return x.StartsWith(y)
	case ew:
		return x.EndsWith(y)
	default:
		panic(fmt.Sprintf("invalid comparison operator %s", b.op))
	}
}

func (b *binary) String() string {
	return fmt.Sprintf("%v %s %v", b.x, b.op, b.y)
}

type logical struct {
	x, y Expr
	op   string
}

func (l *logical) Compare(e Env) bool {
	x := l.x
	y := l.y
	switch l.op {
	case all:
		return x.Compare(e) && y.Compare(e)
	case any:
		return x.Compare(e) || y.Compare(e)
	default:
		panic(fmt.Sprintf("invalid logical operator %s", l.op))
	}
}

func (l *logical) String() string {
	return fmt.Sprintf("%v %s %v", l.x, l.op, l.y)
}
