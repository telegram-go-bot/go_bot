package domain

import "time"

// ChatInfo - describes chat room info known to bot
type ChatInfo struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	CreatedAt time.Time
	Name      string
	ChatID    int64 `gorm:"not null;unique"`
	Enabled   bool  `gorm:"not null;default:'true'"`
}

// ChatUserInfo - detailed user info
type ChatUserInfo struct {
	ID           uint `gorm:"primary_key"`
	SentMessages uint
	FirstName    string
	LastName     string
}

// ChatUser - describes one chat member
type ChatUser struct {
	ID           uint `gorm:"AUTO_INCREMENT"`
	CreatedAt    time.Time
	UserID       int    `gorm:"unique"`
	UserName     string `gorm:"unique"`
	ChatID       int64
	ChatUserInfo ChatUserInfo `gorm:"foreignkey:ID"`
}

// BasicTables - list of known basic tables
var BasicTables = []interface{}{ChatInfo{}, ChatUserInfo{}, ChatUser{}}
