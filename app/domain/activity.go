package domain

// Activity - received notification
// Telegram receives message and this one describes it
type Activity struct {
	MesssageID     int
	Text           string
	ChatID         int64
	LeftChatMember *User
	NewChatMembers *[]User
	RepliedTo      *Activity
	Command        string // comamnd, sent straight to bot using messenger
}

// User - user representation
type User struct {
	ID       int
	UserName string
	IsBot    bool
}
