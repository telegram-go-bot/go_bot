package output

// ShowMessageData - feed presenter's ShowMessage with this
type ShowMessageData struct {
	ChatID       int64
	ReplyToMsgID int
	Text         string
}

// ShowImageData - feed presenter's ShowImage with this
type ShowImageData struct {
	ImageURL string
	ShowMessageData
}

// IPresenter - prepares data to display it via View
// return: @sent_message_id, error
type IPresenter interface {
	ShowMessage(data ShowMessageData) (int, error)
	ShowImage(data ShowImageData) (int, error)
}
