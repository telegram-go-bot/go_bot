package cmn

import "testing"

func TestGetOneMsgFromMany(t *testing.T) {
	msg := getOneMsgFromMany([]string{})
	if len(msg) != 0 {
		t.Error("GetOneMsgFromMany failed on empty array")
	}
}
