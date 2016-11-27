package expr

import (
	"testing"
)

func TestParseWithCast(t *testing.T) {
	e := New()
	e.Varchar("s0", "1.0")

	tests := []struct {
		source string
		want   bool
	}{
		{source: `s0:integer >= 1.0`, want: true},
	}

	for ix, test := range tests {
		ix++
		expr, err := Parse(test.source)
		if err != nil {
			t.Errorf("#%03d: unexpected error: %s", ix, err)
			continue
		}
		if expr.Compare(e) != test.want {
			t.Errorf("#%03d: %s != %t", ix, expr, test.want)
		}
	}
}

func TestParse(t *testing.T) {
	e := New()
	e.Varchar("v0", "world")
	e.Number("n0", 1.0)

	tests := []struct {
		source string
		want   bool
	}{
		{source: `v0 == "world"`, want: true},
		{source: `v0 != "world"`, want: false},
		{source: `n0 == 1.0`, want: true},
		{source: `n0 != 1.0`, want: false},
		{source: `n0 != 2.0`, want: true},
		{source: `v0 == "world" && n0 == 1.0`, want: true},
		{source: `v0 != "world" || n0 == 1.0`, want: true},
		{source: `v0 == "world" || n0 != 1.0`, want: true},
		{source: `v0 == "world" && n0 >= 1.0`, want: true},
		{source: `v0 == "world" && n0 >= 0.0 && n0 <= 2.0`, want: true},
		{source: `v0 >= 1.0`, want: false},
	}

	for ix, test := range tests {
		ix++
		expr, err := Parse(test.source)
		if err != nil {
			t.Errorf("#%03d: unexpected error: %s", ix, err)
			continue
		}
		if expr.Compare(e) != test.want {
			t.Errorf("#%03d: %s != %t", ix, expr, test.want)
		}
	}
}
