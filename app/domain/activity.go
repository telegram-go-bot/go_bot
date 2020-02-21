package domain

// Activity - received notification
// Telegram receives message and this one describes it
type Activity struct {
	Text   string
	ChatID int64
}
