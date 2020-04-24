package olympiadnik

import (
	"strings"

	cmn "github.com/telegram-go-bot/go_bot/app/common"
	websearch "github.com/telegram-go-bot/go_bot/app/common/web_search"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

// Impl - implementation
type Impl struct {
	presenter output.IPresenter
	searcher  websearch.Searcher
}

// New - constructor
func New(presenter output.IPresenter, searcher websearch.Searcher) Impl {
	var tmp = Impl{presenter: presenter}
	tmp.searcher = searcher
	return tmp
}

// OnHelp - display help
func (p Impl) OnHelp() string {
	return ""
}

// OnCommand -
func (p Impl) OnCommand(item raw.Activity) (bool, error) {

	if strings.Index(strings.ToLower(item.Text), "олимпиадник") == -1 {
		return false, nil
	}

	kirill := "кирилл каймаков"

	strToFind := cmn.GetOneMsgFromMany(
		kirill, kirill, kirill, kirill, kirill,
		"победитель олимпиады по физике",
		"победитель олимпиады по программированию",
		"победитель олимпиады по математике",
		kirill, kirill, kirill,
		"николай дуров")

	images := p.searcher.SearchImage(strToFind, 10)
	if len(images) == 0 {
		return false, nil
	}

	pickN := cmn.Rnd.Intn(len(images))
	imageURL := images[pickN]

	if len(imageURL) == 0 {
		return false, nil
	}

	p.presenter.ShowImage(output.ShowImageData{
		ImageURL:        imageURL,
		ShowMessageData: output.ShowMessageData{ChatID: item.ChatID}})

	return false, nil // hidden capture
}
