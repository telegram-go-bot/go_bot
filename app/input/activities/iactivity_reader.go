package activities

import (
	raw "github.com/telegram-go-bot/go_bot/app/domain"
)

//go:generate mockgen -destination=../mocks/mock_iactivity_reader.go -package=mocks github.com/telegram-go-bot/go_bot/app/input/activities IActivityReader

// IActivityReader -
type IActivityReader interface {
	GetActivity() (raw.Activity, error)
}
