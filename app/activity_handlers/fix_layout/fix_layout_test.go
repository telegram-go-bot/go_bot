package fixlayout

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/mocks"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

func executeFixtextTest(t *testing.T, message string, expectedString string) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPresenter := mocks.NewMockIPresenter(mockCtrl)
	commandHandler := New(mockPresenter)

	parentActivity := raw.Activity{ChatID: 111, Text: message}
	testActivity := raw.Activity{ChatID: 111, Text: "!fix", RepliedTo: &parentActivity}

	mockPresenter.EXPECT().ShowMessage(output.ShowMessageData{ChatID: 111, Text: expectedString}).MaxTimes(1)

	succ, err := commandHandler.OnCommand(testActivity)
	if err != nil {
		t.Error(err.Error())
	}
	if !succ {
		t.Error(t.Name() + " failed")
	}
}

func TestDetectLanguage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPresenter := mocks.NewMockIPresenter(mockCtrl)
	New(mockPresenter)

	rus, eng := delectLang("english")
	if rus != false || eng != true {
		t.Error("English lang detection broken")
	}

	rus, eng = delectLang("рус")
	if rus != true || eng != false {
		t.Error("Rus lang detection broken")
	}

	rus, eng = delectLang("enру")
	if rus != false || eng != true {
		t.Error("Rus+Eng detection broken")
	}

	// none
	rus, eng = delectLang("-=-,!*()12345")
	if rus != false || eng != false {
		t.Error("None lang detection broken")
	}
}

func TestBasic(t *testing.T) {
	executeFixtextTest(t, "ghbdtn", "привет")
}

func TestBasic2(t *testing.T) {
	executeFixtextTest(t, "привет", "ghbdtn")
}

func TestBasicSmall(t *testing.T) {
	executeFixtextTest(t, "П", "G")
}
