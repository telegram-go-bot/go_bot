package activityhandlers

import raw "github.com/telegram-go-bot/go_bot/app/domain"

//go:generate mockgen -destination=../mocks/mock_icommand_handler.go -package=mocks github.com/telegram-go-bot/go_bot/app/activity_handlers ICommandHandler

// ICommandHandler - common interface for all handlers
type ICommandHandler interface {
	OnHelp() string
	OnCommand(raw.Activity) (bool, error)
}
