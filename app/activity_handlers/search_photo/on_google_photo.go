package googlephoto

import (
	"errors"
	"log"
	"net/url"

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

/*
	// if one query was repeated several times - increase set of images to choose from
	if strToFind == prevQuery {
		prevDepthVal = prevDepthVal + 5
	} else {
		prevQuery = strToFind
		prevDepthVal = google.GMaxImagesResult
	}

	return OnGooglePhotoImpl(chatID, strToFind, prevDepthVal, bot)
}

//------------------------------------------------------------------------------------------------------
*/
// OnGooglePhotoImpl - underlying (private) impl of image querier
func (p impl) onGooglePhotoImpl(item raw.Activity, strToFind string, depth int) (bool, error) {

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	var imageURL string
	var images []string
	// search for query string for an image
	_, err := url.ParseRequestURI(strToFind)
	if err != nil {
		images = p.searcher.SearchImage(strToFind, depth)
		if len(images) == 0 {
			SendMsg("Missing Data...")
			return false, errors.New("0 images found querying for a text : \"" + strToFind + "\"")
		}
		pickN := cmn.Rnd.Intn(len(images))
		imageURL = images[pickN]
	} else { // if parameter is url - use it
		imageURL = strToFind
	}

	if len(imageURL) == 0 {
		SendMsg(cmn.GetFailMsg())
		return false, errors.New("Empty image URL found querying for a text : \"" + strToFind + "\"")
	}

	log.Printf("Selected image url: " + imageURL)

	p.presenter.ShowImage(output.ShowImageData{
		ImageURL:        imageURL,
		ShowMessageData: output.ShowMessageData{ChatID: item.ChatID}})

	return true, nil
}
