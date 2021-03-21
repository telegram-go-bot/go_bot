package text_to_speech

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	"github.com/telegram-go-bot/go_bot/app/common/vk"
	"github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/output"
)

var mp3UrlRe = regexp.MustCompile(`\(mp3\): (.*?)\n`)

type impl struct {
	presenter output.IPresenter
	lock      *sync.Mutex
}

// New - constructor
func New(presenter output.IPresenter) impl {
	var tmp = impl{presenter: presenter, lock: &sync.Mutex{}}
	return tmp
}

// OnHelp - display help
func (p impl) OnHelp() string {
	return "<b>!голосом|голос</b> - text to speech ivona"
}

// OnCommand -
func (p impl) OnCommand(item domain.Activity) (bool, error) {

	msg, isThisCommand := helpers.IsOnCommand(item.Text, []string{"голосом", "голос", "гг"})
	if !isThisCommand {
		return false, nil
	}

	var text string
	if len(msg) != 0 {
		text = msg
	} else if item.RepliedTo != nil {
		text = item.RepliedTo.Text
	}

	if len(text) == 0 {
		return false, nil
	}

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	p.lock.Lock()
	speechUrl, err := textToSpeech(text)
	p.lock.Unlock()
	if err != nil {
		SendMsg(cmn.GetFailMsg())
		return true, err
	}

	replyToID := item.MesssageID
	if item.RepliedTo != nil {
		replyToID = item.RepliedTo.MesssageID
	}

	return p.showAudioByURL(speechUrl, item.ChatID, replyToID)
}

func (p impl) showAudioByURL(url string, chatID int64, replyToMsgID int) (bool, error) {
	if len(url) == 0 {
		p.presenter.ShowMessage(output.ShowMessageData{ChatID: chatID, Text: cmn.GetFailMsg()})
		return false, errors.New("Empty audio URL found querying for a text-URL")
	}

	log.Printf("showAudioByURL: " + url)

	_, err := p.presenter.ShowAudio(output.ShowAudioData{
		AudioURL: url,
		ShowMessageData: output.ShowMessageData{
			ChatID:       chatID,
			ReplyToMsgID: replyToMsgID}})
	if err != nil {
		return false, err
	}
	return true, nil
}

//----------------------------------------------------------------------------------------------------------

func getUrl(msg string) (string, bool) {
	match := mp3UrlRe.FindStringSubmatch(msg)
	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}

func waitForResponse(convID string) (string, error) {
	parameters := make(map[string]string)

	userID := os.Getenv("IVONA_TEXT_TO_SPEECH_USER_ID")
	if len(convID) == 0 {
		return "", fmt.Errorf("IVONA_TEXT_TO_SPEECH_USER_ID env var is EMPTY!!!")
	}

	parameters["user_id"] = userID
	parameters["count"] = "1"
	parameters["peer_id"] = convID

	resp, err := vk.Request("messages.getHistory", parameters)
	if err != nil {
		return "", fmt.Errorf("messages.getHistory request failed: %s", err)
	}

	fmt.Println(string(resp))

	var vkResponse vk.MessagesGetHistoryResponse

	err = json.Unmarshal(resp, &vkResponse)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling json response: %s", err)
	}

	if vkResponse.Response.Count == 0 {
		return "", nil
	}

	url, parsed := getUrl(vkResponse.Response.Items[0].Text)
	if !parsed {
		return "", fmt.Errorf("Error extracting url from: %s", vkResponse.Response.Items[0].Text)
	}

	return url, nil
}

// returns url to audio message
func textToSpeech(msg string) (string, error) {

	convID := os.Getenv("IVONA_TEXT_TO_SPEECH_CONVERSATION_ID")
	if len(convID) == 0 {
		return "", fmt.Errorf("IVONA_TEXT_TO_SPEECH_CONVERSATION_ID env var is EMPTY!!!")
	}
	parameters := make(map[string]string)
	parameters["message"] = msg
	parameters["peer_id"] = convID
	parameters["random_id"] = strconv.Itoa(cmn.Rnd.Intn(1000000))

	resp, err := vk.Request("messages.send", parameters)
	if err != nil {
		return "", fmt.Errorf("messages.send request failed: %s", err)
	}

	fmt.Println(string(resp))

	for {
		time.Sleep(1 * time.Second)
		msg, err := waitForResponse(convID)
		if err != nil {
			return "", err
		}

		if len(msg) != 0 {
			return msg, nil
		}
	}
}
