package activityhandlers

import raw "github.com/telegram-go-bot/go_bot/app/domain"

// ICommandHandler - common interface for all handlers
type ICommandHandler interface {
	OnHelp() string
	OnCommand(raw.Activity) (bool, error)
}
