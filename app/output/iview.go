package output

// IView -
type IView interface {
	ShowMessage(msg ViewMessageData) error
}
