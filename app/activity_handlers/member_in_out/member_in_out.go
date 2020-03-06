package memberinout

import (
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

type impl struct {
	presenter output.IPresenter
}

// New - constructor
func New(presenter output.IPresenter) impl {
	return impl{presenter: presenter}
}

// OnHelp - display help
func (p impl) OnHelp() string {
	// NO help! Hidden cmd ;D
	return ""
}

func (p impl) OnCommand(item raw.Activity) (bool, error) {

	if item.LeftChatMember != nil {
		err := p.OnMemberLeft(item.LeftChatMember)
		return true, err
	}

	if item.NewChatMembers != nil {
		err := p.OnMemberJoined(*item.NewChatMembers)
		return true, err
	}

	return false, nil
}

func (p impl) OnMemberLeft(user *raw.User) error {
	return nil
}

func (p impl) OnMemberJoined(users []raw.User) error {
	return nil
}
