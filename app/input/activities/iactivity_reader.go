package activities

import (
	raw "github.com/telegram-go-bot/go_bot/app/domain"
)

// IActivityReader -
type IActivityReader interface {
	GetActivity() (raw.Activity, error)
}
