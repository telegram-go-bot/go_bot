package googlephoto

import "testing"

func TestErasion(t *testing.T) {
	a := []string{"A", "B", "C", "D", "E"}

	r := erase(a, 0)
	if r[0] == "A" {
		t.Error("Erase from begining failed")
	}

	r = erase(a, 4)
	if r[3] != "D" {
		t.Error("Erase from the end failed")
	}

	r = erase(r, 1)
	if r[1] != "D" {
		t.Error("Erase from the end failed")
	}
}
