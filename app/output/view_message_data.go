package output

// ViewMessageData - simple message data to be shown by IView
type ViewMessageData struct {
	ReplyToMsgID int
	ChatID       int64
	Text         string
	ParseMode    int
}

// ViewImageData - shows image
type ViewImageData struct {
	ImageData []byte
	ViewMessageData
}

// ViewAnimationData - shows image
type ViewAnimationData struct {
	AnimationID string
	Caption     string
	ViewMessageData
}

// ViewAudioData - uploads audio and shows it
type ViewAudioData struct {
	AudioData []byte
	Caption   string
	ViewMessageData
}
