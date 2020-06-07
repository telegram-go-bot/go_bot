package settings

import (
	raw "github.com/telegram-go-bot/go_bot/app/domain"
)

// ISettingsOpener - access settings
type ISettingsOpener interface {
	Init(url string) error
}

// ISettingsReader - readonly access
type ISettingsReader interface {
	GetChatInfo(chatID int64) (*raw.ChatInfo, error)
	GetChatUser(userID int) (*raw.ChatUser, error)
	// get all db entries, related to a definite handler (e.g. settings for a Handler)
	// Mapping interface{} to struct is up to handler
	GetHandlerRecords(out interface{}) error
}

// ISettingsWriter - add new entries
type ISettingsWriter interface {
	AddRecord(newRec interface{}) error
	UpdateRecord(modifiedRec interface{}) error
}

// ISettings - repository
type ISettings interface {
	ISettingsOpener
	ISettingsReader
	ISettingsWriter
}

var (
	// Inst - settings singleton
	inst ISettings
)

// New - init new
func New(settings ISettings) {
	inst = settings
}

// Inst -singletone getter. returns already created instance
func Inst() ISettings {
	return inst
}
