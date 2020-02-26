package output

// ViewMessageData - simple message data to be shown by IView
type ViewMessageData struct {
	ReplyToMsgID int
	ChatID       int64
	Text         string
}

// ViewImageData - shows image
type ViewImageData struct {
	ImageData []byte
	ViewMessageData
}
