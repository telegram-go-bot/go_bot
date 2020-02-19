package output

// ShowMessageData - feed presenter's ShowMessage with this
type ShowMessageData struct {
	ChatID       int64
	ReplyToMsgID int
	Text         string
}

// IPresenter - prepares data to display it via View
// return: @sent_message_id, error
type IPresenter interface {
	ShowMessage(data ShowMessageData) (int, error)
}
