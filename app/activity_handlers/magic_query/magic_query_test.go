package magicquery

import "testing"

func TestParseAddBase(t *testing.T) {
	magic, val, typ, err := parseAdd("Magic" + queryDelimiter + "Val" + queryDelimiter + "Type")
	if err != nil ||
		magic != "magic" ||
		typ != "Type" ||
		val != "val" {
		t.Error("unexpected add_query param format")
	}

	_, _, _, err = parseAdd("magic")
	if err == nil {
		t.Error("unexpected add_query param format")
	}
}

func TestParseAddFriendly(t *testing.T) {
	magic, val, typ, err := parseAdd("Magic" + queryDelimiter + "val")
	if err != nil ||
		magic != "magic" ||
		typ != defaultValueType ||
		val != "val" {
		t.Error("expected default type")
	}
}
