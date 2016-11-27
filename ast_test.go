package expr

import (
	"testing"
)

func TestBinaryCompareNumber(t *testing.T) {
	tests := []struct {
		x, y number
		op   string
		want bool
	}{
		{x: number(1.0), y: number(1.0), op: eq, want: true},
		{x: number(1.0), y: number(1.0), op: ne, want: false},
		{x: number(1.0), y: number(2.0), op: eq, want: false},
		{x: number(1.0), y: number(2.0), op: ne, want: true},
		{x: number(1.0), y: number(2.0), op: lt, want: true},
		{x: number(1.0), y: number(2.0), op: le, want: true},
		{x: number(2.0), y: number(1.0), op: gt, want: true},
		{x: number(2.0), y: number(1.0), op: ge, want: true},
	}

	for ix, test := range tests {
		b := &binary{test.x, test.y, test.op}
		if b.Compare(nil) != test.want {
			t.Errorf("#%d: %q failed => expected: %t", ix, b, test.want)
		}
	}
}

func TestBinaryString(t *testing.T) {
	tests := []struct {
		x, y Value
		op   string
		want string
	}{
		{x: varchar("hello"), y: varchar("world"), op: eq, want: "hello == world"},
		{x: varchar("hello"), y: varchar("world"), op: ne, want: "hello != world"},
	}

	for ix, test := range tests {
		b := &binary{x: test.x, y: test.y, op: test.op}
		if s := b.String(); s != test.want {
			t.Errorf("#%d) got %s, want %s", ix, s, test.want)
		}
	}
}

func TestBinaryCompareVarchar(t *testing.T) {
	v1 := varchar("hello")
	v2 := varchar("world")
	v3 := varchar("abcd")
	v4 := varchar("efgh")

	tests := []struct {
		b    Expr
		want bool
	}{
		{&binary{x: v1, y: v2, op: eq}, false},
		{&binary{x: v1, y: v2, op: ne}, true},
		{&binary{x: v1, y: v2, op: ne}, true},
		{&binary{x: v3, y: v4, op: lt}, true},
		{&binary{x: v3, y: v4, op: le}, true},
		{&binary{x: v4, y: v3, op: gt}, true},
		{&binary{x: v4, y: v3, op: ge}, true},
	}

	for ix, test := range tests {
		if test.b.Compare(nil) != test.want {
			t.Errorf("#%d: %q failed => expected: %t", ix, test.b, test.want)
		}
	}
}

func TestLogicalCompareSimple(t *testing.T) {
	tests := []struct {
		x, y Expr
		op   string
		want bool
	}{
		{
			x:    &binary{varchar("hello"), varchar("hello"), eq},
			y:    &binary{number(1.0), number(1.0), eq},
			op:   all,
			want: true,
		},
		{
			x:    &binary{varchar("hello"), varchar("hello"), eq},
			y:    &binary{number(1.0), number(1.0), eq},
			op:   any,
			want: true,
		},
	}

	for ix, test := range tests {
		l := &logical{test.x, test.y, test.op}
		if l.Compare(nil) != test.want {
			t.Errorf("#%d: %q failed => expected => %t", ix, l, test.want)
		}
	}
}
