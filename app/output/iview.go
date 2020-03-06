package output

// IView - display smth
// return: @sent_message_id, error
type IView interface {
	ShowMessage(msg ViewMessageData) (int, error)
	ShowImage(msg ViewImageData) (int, error)
	ShowAnimation(msg ViewAnimationData) (int, error)
}
