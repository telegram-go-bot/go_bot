package output

// ViewMessageData - simple message data to be shown by IView
type ViewMessageData struct {
	ReplyToMsgID int64
	ChatID       int64
	Text         string
}
