package activityhelpers

import (
	"testing"
)

func TestIsOnCommandNoCmd(t *testing.T) {
	res, ok := IsOnCommand("cmd", []string{"cmd"})
	if ok || res != "" {
		t.Error("No command passed by some detected")
	}
}

func TestIsOnCommand(t *testing.T) {
	_, ok := IsOnCommand("!cmd", []string{"cmd"})
	if !ok {
		t.Error("IsOnCommand basic invalid")
	}

	_, ok = IsOnCommand("!cmd", []string{"cmdcmd"})
	if ok {
		t.Error("IsOnCommand including cmd fail")
	}

	_, ok = IsOnCommand("!chicken", []string{"aeiouy"})
	if ok {
		t.Error("IsOnCommand invalid cmd fail")
	}

	cmd, ok := IsOnCommand("noo !cmd <3", []string{"cmd"})
	if !ok || cmd != "<3" {
		t.Error("IsOnCommand failed - cmd might be not at the beginning of the message")
	}
}

func TestIsOnCommandParams(t *testing.T) {
	params, ok := IsOnCommand("!cmd param1 param2", []string{"cmd"})
	if !ok {
		t.Error("IsOnCommand with params failed parsing cmd")
	}

	if params != "param1 param2" {
		t.Error("IsOnCommand with params failed getting params")
	}
}
