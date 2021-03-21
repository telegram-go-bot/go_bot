package output

const (
	// ParseModeMarkdown - simple ** __ etc. formatting
	ParseModeMarkdown = iota
	// ParseModeHTML - treat message as html (inc formatiing etc.)
	ParseModeHTML
)

// ShowMessageData - feed presenter's ShowMessage with this
type ShowMessageData struct {
	ChatID       int64
	ReplyToMsgID int
	Text         string
	ParseMode    int
}

// ShowImageData - feed presenter's ShowImage with this
type ShowImageData struct {
	ImageURL     string
	RawImageData []byte
	ShowMessageData
}

// ShowAnimationData - feed presenter's ShowImage with this
type ShowAnimationData struct {
	AnimationID string
	Caption     string
	ShowMessageData
}

// ShowAudioData - feed presenter's ShowAudio with this
type ShowAudioData struct {
	AudioURL     string
	RawAudioData []byte
	Caption      string
	ShowMessageData
}

//go:generate mockgen -destination=../mocks/mock_ipresenter.go -package=mocks github.com/telegram-go-bot/go_bot/app/output IPresenter

// IPresenter - prepares data to display it via View
// return: @sent_message_id, error
type IPresenter interface {
	ShowMessage(data ShowMessageData) (int, error)
	ShowImage(data ShowImageData) (int, error)
	ShowGif(data ShowAnimationData) (int, error)
	ShowAudio(data ShowAudioData) (int, error)
}
