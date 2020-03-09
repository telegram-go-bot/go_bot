package cmn

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestGetOneMsgFromMany(t *testing.T) {
	msg := GetOneMsgFromMany()
	if len(msg) != 0 {
		t.Error("GetOneMsgFromMany failed on empty array")
	}
}

func TestGetOneMsgFromOne(t *testing.T) {
	msg := GetOneMsgFromMany("test")
	if msg != "test" {
		t.Error("GetOneMsgFromOne failed")
	}
}

func TestGetFailMsgReturnsSmth(t *testing.T) {
	msg := GetFailMsg()
	if len(msg) == 0 {
		t.Error("GetFailMsg failed returning empty string")
	}
}
