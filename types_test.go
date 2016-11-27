package expr

import (
	"testing"
)

func TestIncompatibleType(t *testing.T) {
	v := varchar("hello")
	n := number(1.0)

	if v.Equal(n) {
		t.Errorf("v == n => impossible")
	}
}

func TestNumberEq(t *testing.T) {
	n1 := number(0.0)
	n2 := number(1.0)

	if n1.Equal(n2) {
		t.Errorf("%f == %f", n1, n2)
	}
	if !n1.Equal(n1) {
		t.Errorf("%f != %f", n1, n1)
	}
}

func TestNumberLess(t *testing.T) {
	n1 := number(0.0)
	n2 := number(1.0)

	if !n1.Less(n2) {
		t.Errorf("%f > %f", n2, n1)
	}
	if n2.Less(n1) {
		t.Errorf("%f < %f", n2, n1)
	}
}

func TestVarcharEq(t *testing.T) {
	s1 := varchar("hello")
	s2 := varchar("world")

	if s1.Equal(s2) {
		t.Errorf("%s == %s", s1, s2)
	}
	if !s1.Equal(s1) {
		t.Errorf("%s != %s", s1, s1)
	}
}

func TestVarcharLess(t *testing.T) {
	s1 := varchar("abcd")
	s2 := varchar("efgh")
	if s2.Less(s1) {
		t.Errorf("%s < %s", s2, s1)
	}
	if !s1.Less(s2) {
		t.Errorf("%s > %s", s1, s2)
	}
}
