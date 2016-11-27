package expr

import (
	"errors"
	"fmt"
	"time"
)

var (
	MalformedErr     = errors.New("malformed input string")
	UnsupportedOpErr = errors.New("invalid operator")
	NotFoundErr      = errors.New("not a value")
)

var invalid = empty{}

const (
	eq = "=="
	ne = "!="
	ge = ">="
	le = "<="
	gt = ">"
	lt = "<"
	al = "~="
	sw = "^="
	ew = "$="
)

const (
	all = "&&"
	any = "||"
)

type Type int

const (
	Null Type = iota
	Default
	Varchar
	Number
	Integer
	Real
	Datetime
)

var dtypes = map[string]Type{
	"varchar":  Varchar,
	"string":   Varchar,
	"real":     Number,
	"number":   Number,
	"datetime": Datetime,
	"moment":   Datetime,
}

type Expr interface {
	Compare(Env) bool
	fmt.Stringer
}

type Value interface {
	Equal(Value) bool
	Less(Value) bool
	Almost(Value) bool
	StartsWith(Value) bool
	EndsWith(Value) bool
	Type() Type
	Lookup(Env) Value
	Cast(Type) Value
}

type Env map[string]Value

func New() Env {
	return Env{}
}

func (e Env) Varchar(name string, value string) {
	e[name] = varchar(value)
}

func (e Env) Number(name string, value float64) {
	e[name] = number(value)
}

func (e Env) Datetime(name string, value time.Time) {
	e[name] = datetime(value)
}

func (e Env) lookup(name string) Value {
	if e == nil || len(e) == 0 {
		return invalid
	}
	if v, ok := e[name]; ok {
		return v
	} else {
		return invalid
	}
}

type Var struct {
	name  string
	dtype Type
	value Value
}

func (v Var) Type() Type {
	return v.value.Type()
}

func (v Var) Equal(other Value) bool {
	return v.value.Equal(other)
}

func (v Var) Less(other Value) bool {
	return v.value.Less(other)
}

func (v Var) Almost(other Value) bool {
	return v.value.Almost(other)
}

func (v Var) StartsWith(other Value) bool {
	return v.value.StartsWith(other)
}

func (v Var) EndsWith(other Value) bool {
	return v.value.EndsWith(other)
}

func (v Var) Lookup(e Env) Value {
	val := e.lookup(v.name)
	if v.dtype != Default {
		val = val.Cast(v.dtype)
	}
	v.value = val
	return v.value
}

func (v Var) String() string {
	return v.name
}

func (v Var) Cast(t Type) Value {
	return v.value.Cast(t)
}
