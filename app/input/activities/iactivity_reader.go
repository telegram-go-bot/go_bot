package activities

import (
	raw "github.com/telegram-go-bot/go_bot/app/domain/raw_structures"
)

// IActivityReader -
type IActivityReader interface {
	GetActivity() (raw.Activity, error)
}
