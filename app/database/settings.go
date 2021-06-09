package settings

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // have to
	raw "github.com/telegram-go-bot/go_bot/app/domain"
)

// Settings - one settings manager for gorm postgres db. ISettings impl
type Settings struct {
	db *gorm.DB
}

// Init - init db connection.
// url - dabase url
func (p *Settings) Init(url string) error {
	var err error
	p.db, err = gorm.Open("postgres", url)
	if err != nil {
		return err
	}
	p.db.LogMode(true)
	p.db.Callback().Create().Remove("gorm:force_reload_after_create") // dont call SELECT after INSERT

	for _, rawType := range raw.BasicTables {
		p.db = p.db.AutoMigrate(rawType)
	}

	return nil
}

// GetChatInfo -
func (p *Settings) GetChatInfo(chatID int64) (*raw.ChatInfo, error) {

	chatRoomInfo := raw.ChatInfo{}
	err := p.db.Where(&raw.ChatInfo{ChatID: chatID}).First(&chatRoomInfo).Error
	if err != nil {
		return nil, err
	}

	return &chatRoomInfo, nil
}

// GetChatUser -
func (p *Settings) GetChatUser(userID int) (*raw.ChatUser, error) {

	chatUser := raw.ChatUser{}
	err := p.db.Where(&raw.ChatUser{UserID: userID}).First(&chatUser).Error
	if err != nil {
		return nil, err
	}

	err = p.db.Model(&chatUser).Related(&chatUser.ChatUserInfo, "ChatUserInfo").Error
	if err != nil {
		return nil, err
	}

	return &chatUser, nil
}

// GetHandlerRecords -
func (p *Settings) GetHandlerRecords(out interface{}) error {
	err := p.db.Find(out).Error
	if err != nil {
		if p.db.HasTable(out) {
			return err
		}
	}
	return nil
}

// AddRecord - insert record to db
func (p *Settings) AddRecord(newRec interface{}) error {

	if !p.db.HasTable(newRec) {
		p.db.CreateTable(newRec)
	}

	err := p.db.Create(newRec).Error
	return err
}

// UpdateRecord - update existing
func (p *Settings) UpdateRecord(newRec interface{}) error {
	// todo: update only defined fields...
	err := p.db.Update(newRec).Error
	return err
}
