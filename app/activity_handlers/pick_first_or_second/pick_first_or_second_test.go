package pickfirstorsecond

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/mocks"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

type Random0 struct {
}

func (r Random0) Intn(int) int {
	return 0
}

func init() {
	cmn.Rnd = Random0{}
}

func TestBasic(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPresenter := mocks.NewMockIPresenter(mockCtrl)
	commandHandler := New(mockPresenter)

	testActivity := raw.Activity{ChatID: 111, Text: "!billy нет или да"}

	mockPresenter.EXPECT().ShowMessage(output.ShowMessageData{ChatID: 111, Text: "да"}).MaxTimes(1)

	succ, err := commandHandler.OnCommand(testActivity)
	if err != nil {
		t.Error(err.Error())
	}
	if !succ {
		t.Error("BasicTest failed")
	}
}
