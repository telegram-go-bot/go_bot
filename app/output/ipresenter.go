package output

// ShowMessageData - feed presenter's ShowMessage with this
type ShowMessageData struct {
	ChatID int64
	Text   string
}

// IPresenter - prepares data to display it via View
type IPresenter interface {
	ShowMessage(data ShowMessageData) error
}
