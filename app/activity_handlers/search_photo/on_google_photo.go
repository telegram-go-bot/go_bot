package googlephoto

import (
	"errors"
	"log"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	websearch "github.com/telegram-go-bot/go_bot/app/common/web_search"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/output"
)

/*/*****************************************************************
 TO DO:
 1. query for MORE that 10 images ( is not working yet )
/*/

type impl struct {
	presenter    output.IPresenter
	prevQuery    string
	prevDepthVal int
	searcher     websearch.Searcher
}

// New - constructor
func New(presenter output.IPresenter, searcher websearch.Searcher) impl {
	var tmp = impl{presenter: presenter}
	tmp.prevDepthVal = searcher.MaxResponseItemsCount()
	tmp.searcher = searcher
	return tmp
}

// OnHelp - display help
func (p impl) OnHelp() string {
	return "<b>!фото|photo</b> <i>стринга</i> - рандомная фотка из 10 первых результатов запроса"
}

// OnCommand -
func (p impl) OnCommand(item raw.Activity) (bool, error) {

	strToFind, isThisCommand := helpers.IsOnCommand(item.Text, []string{"фото", "photo"})
	if !isThisCommand {
		return false, nil
	}

	return p.onGooglePhotoImpl(item, strToFind, p.prevDepthVal)
}

func (p impl) showImageByURL(url string, chatID int64) (bool, error) {
	if len(url) == 0 {
		p.presenter.ShowMessage(output.ShowMessageData{ChatID: chatID, Text: cmn.GetFailMsg()})
		return false, errors.New("Empty image URL found querying for a text-URL")
	}

	log.Printf("ShowImageByURL: " + url)

	_, err := p.presenter.ShowImage(output.ShowImageData{
		ImageURL:        url,
		ShowMessageData: output.ShowMessageData{ChatID: chatID}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func erase(a []string, idx int) []string {
	a[idx] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = ""     // Erase last element (write zero value).
	a = a[:len(a)-1]     // Truncate slice.
	return a
}

// OnGooglePhotoImpl - underlying (private) impl of image querier
func (p impl) onGooglePhotoImpl(item raw.Activity, strToFind string, depth int) (bool, error) {

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	var imageURL string
	var images []string

	images = p.searcher.SearchImage(strToFind, depth)
	if len(images) == 0 {
		SendMsg("Missing Data...")
		return true, errors.New("0 images found querying for a text : \"" + strToFind + "\"")
	}

	for attempts := len(images); attempts != 0; attempts-- {

		if len(images) == 0 {
			SendMsg(cmn.GetFailMsg())
			return true, errors.New("No more URLS left for a text : \"" + strToFind + "\"")
		}

		pickN := cmn.Rnd.Intn(len(images))
		imageURL = images[pickN]

		images = erase(images, pickN)

		done, err := p.showImageByURL(imageURL, item.ChatID)
		if done && err == nil {
			return true, nil
		}

		continue
	}

	return true, errors.New("No more URLS left for a text : \"" + strToFind + "\"")
}
