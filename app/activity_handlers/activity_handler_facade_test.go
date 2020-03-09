package activityhandlers

import (
	"fmt"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/mocks"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

func TestNewFacade(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPresenter := mocks.NewMockIPresenter(mockCtrl)
	commandHandler := mocks.NewMockICommandHandler(mockCtrl)

	handlers := []ICommandHandler{commandHandler}

	const dummyHelpStr = "some help msg"
	commandHandler.EXPECT().OnHelp().Return(dummyHelpStr).Times(1)

	newFacade := New(handlers, mockPresenter)
	if !strings.Contains(newFacade.helpMsg, dummyHelpStr) {
		t.Error("Error generating help string")
	}
	if len(newFacade.handlers) != 1 {
		t.Error("Too much command handlers")
	}
	if newFacade.handlers[0] != commandHandler {
		t.Error("Unexpected command handler")
	}
}

type ShowMessageDataMatcher struct {
	dat output.ShowMessageData
}

func CheckShowMessageData(dat output.ShowMessageData) gomock.Matcher {
	return &ShowMessageDataMatcher{dat: dat}
}

func (o *ShowMessageDataMatcher) Matches(x interface{}) bool {
	right := x.(output.ShowMessageData)
	return o.dat.ChatID == right.ChatID && o.dat.ParseMode == right.ParseMode
}

func (o *ShowMessageDataMatcher) String() string {
	return fmt.Sprintf("%v", o.dat)
}

func TestProcessActivitiesDisplayHelpMsg(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPresenter := mocks.NewMockIPresenter(mockCtrl)
	commandHandler := mocks.NewMockICommandHandler(mockCtrl)

	handlers := []ICommandHandler{commandHandler}
	commandHandler.EXPECT().OnHelp().Return("dummyHelpStr").Times(1)

	newFacade := New(handlers, mockPresenter)

	testActivity := raw.Activity{ChatID: 111}

	mockPresenter.EXPECT().ShowMessage(CheckShowMessageData(
		output.ShowMessageData{
			ChatID:    111,
			ParseMode: output.ParseModeHTML}))

	newFacade.displayHelpMsg(&testActivity)
}
