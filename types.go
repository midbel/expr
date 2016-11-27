package expr

import (
	"math"
	"strconv"
	"strings"
	"time"
)

const IsoFmt = "2006-01-02 15:04:05.999999"

type empty struct{}

func (e empty) Type() Type {
	return Null
}

func (e empty) Equal(_ Value) bool {
	return false
}

func (e empty) Less(_ Value) bool {
	return false
}

func (e empty) Almost(_ Value) bool {
	return false
}

func (e empty) StartsWith(_ Value) bool {
	return false
}

func (e empty) EndsWith(_ Value) bool {
	return false
}

func (e empty) Bool() bool {
	return false
}

func (e empty) String() string {
	return "invalid"
}

func (e empty) Lookup(_ Env) Value {
	return e
}

func (e empty) Cast(_ Type) Value {
	return e
}

type number float64

func (n number) Type() Type {
	return Number
}

func (n number) Equal(other Value) bool {
	if other, ok := other.(number); ok {
		return float64(n) == float64(other)
	}
	return false
}

func (n number) Less(other Value) bool {
	if other, ok := other.(number); ok {
		return float64(n) < float64(other)
	}
	return false
}

func (n number) Bool() bool {
	return float64(n) != 0.0
}

func (n number) Almost(other Value) bool {
	if other, ok := other.(number); ok {
		return math.Floor(float64(n)) == float64(other)
	}
	return false
}

func (n number) StartsWith(other Value) bool {
	return n.Equal(other) || !n.Less(other)
}

func (n number) EndsWith(other Value) bool {
	return n.Equal(other) || n.Less(other)
}

func (n number) Lookup(_ Env) Value {
	return n
}

func (n number) Cast(t Type) Value {
	switch t {
	case Varchar:
		v := strconv.FormatFloat(float64(t), 'g', 6, 64)
		return varchar(v)
	case Number:
		return n
	case Datetime:
		return invalid
	default:
		return invalid
	}
}

type varchar string

func (v varchar) Type() Type {
	return Varchar
}

func (v varchar) Equal(other Value) bool {
	if other, ok := other.(varchar); ok {
		return strings.Compare(string(v), string(other)) == 0
	}
	return false
}

func (v varchar) Less(other Value) bool {
	if other, ok := other.(varchar); ok {
		return strings.Compare(string(v), string(other)) == -1
	}
	return false
}

func (v varchar) Bool() bool {
	return string(v) != ""
}

func (v varchar) Almost(other Value) bool {
	if other, ok := other.(varchar); ok {
		return strings.Contains(string(v), string(other))
	}
	return false
}

func (v varchar) StartsWith(other Value) bool {
	if other, ok := other.(varchar); ok {
		return strings.HasPrefix(string(v), string(other))
	}
	return false
}

func (v varchar) EndsWith(other Value) bool {
	if other, ok := other.(varchar); ok {
		return strings.HasSuffix(string(v), string(other))
	}
	return false
}

func (v varchar) Lookup(_ Env) Value {
	return v
}

func (v varchar) Cast(t Type) Value {
	switch t {
	case Varchar:
		return v
	case Number:
		if v, err := strconv.ParseFloat(string(v), 64); err == nil {
			return number(v)
		} else {
			return invalid
		}
	case Datetime:
		if v, err := time.Parse(IsoFmt, string(v)); err == nil {
			return datetime(v)
		} else {
			return invalid
		}
	default:
		return invalid
	}
}

type datetime time.Time

func (d datetime) Type() Type {
	return Datetime
}

func (d datetime) Equal(other Value) bool {
	if other, ok := other.(datetime); ok {
		return time.Time(d).Equal(time.Time(other))
	}
	return false
}

func (d datetime) Less(other Value) bool {
	if other, ok := other.(datetime); ok {
		return time.Time(d).Before(time.Time(other))
	}
	return false
}

func (d datetime) Almost(other Value) bool {
	return false
}

func (d datetime) StartsWith(other Value) bool {
	return d.Equal(other) || !d.Less(other)
}

func (d datetime) EndsWith(other Value) bool {
	return d.Equal(other) || d.Less(other)
}

func (d datetime) Bool() bool {
	return !time.Time(d).IsZero()
}

func (d datetime) Lookup(_ Env) Value {
	return d
}

func (d datetime) Cast(t Type) Value {
	switch t {
	case Varchar:
		t := time.Time(d)
		return varchar(t.Format(IsoFmt))
	case Number:
		t := time.Time(d)
		return number(float64(t.Unix()))
	case Datetime:
		return d
	default:
		return invalid
	}
}
