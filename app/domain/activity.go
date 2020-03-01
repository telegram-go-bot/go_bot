package domain

// Activity - received notification
// Telegram receives message and this one describes it
type Activity struct {
	Text           string
	ChatID         int64
	LeftChatMember *User
	NewChatmember  *[]User
	RepliedTo      *Activity
}

// User - user representation
type User struct {
	ID       int
	UserName string
	IsBot    bool
}
