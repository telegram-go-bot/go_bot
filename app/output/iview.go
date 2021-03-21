package output

//go:generate mockgen -destination=../mocks/mock_iview.go -package=mocks github.com/telegram-go-bot/go_bot/app/output IView

// IView - display smth
// return: @sent_message_id, error
type IView interface {
	ShowMessage(msg ViewMessageData) (int, error)
	ShowImage(msg ViewImageData) (int, error)
	ShowAnimation(msg ViewAnimationData) (int, error)
	ShowAudio(msg ViewAudioData) (int, error)
}
