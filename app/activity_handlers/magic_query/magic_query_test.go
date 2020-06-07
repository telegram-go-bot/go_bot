package magicquery

import "testing"

func TestParseAddBase(t *testing.T) {
	parsed, err := parseAdd("Magic" + queryDelimiter + "Val" + queryDelimiter + "Type")
	if err != nil ||
		parsed.magic != "magic" ||
		parsed.typ != "Type" ||
		parsed.val != "val" {
		t.Error("unexpected add_query param format")
	}

	_, err = parseAdd("magic")
	if err == nil {
		t.Error("unexpected add_query param format")
	}
}

func TestParseAddFriendly(t *testing.T) {
	parsed, err := parseAdd("Magic" + queryDelimiter + "val")
	if err != nil ||
		parsed.magic != "magic" ||
		parsed.typ != defaultValueType ||
		parsed.val != "val" {
		t.Error("expected default type")
	}
}
