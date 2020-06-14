package covid

import "testing"

func TestToString(t *testing.T) {
	s := toString(123)
	if s != "123" {
		t.Error("expected default type")
	}

	s = toString(1234)
	if s != "1 234" {
		t.Error("unexpected: " + s)
	}

	s = toString(123456)
	if s != "123 456" {
		t.Error("unexpected: " + s)
	}

	s = toString(1234567)
	if s != "1 234 567" {
		t.Error("unexpected: " + s)
	}

	s = toString(12)
	if s != "12" {
		t.Error("unexpected: " + s)
	}

	s = toString(0)
	if s != "0" {
		t.Error("unexpected: " + s)
	}
}

func TestToStringCorner(t *testing.T) {
	s := toString(-1)
	if s != "-1" {
		t.Error("unexpected: " + s)
	}

	s = toString(-123345)
	if s != "-123 345" {
		t.Error("unexpected: " + s)
	}
}
