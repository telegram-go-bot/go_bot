package loopapoopa

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	"github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/output"
)

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
	return "<b>!лупа|пупа|loopa|poopa</b> - анекдоты про лупу и пупу на каждый день"
}

// OnCommand -
func (p impl) OnCommand(item domain.Activity) (bool, error) {

	_, isThisCommand := helpers.IsOnCommand(item.Text, []string{"лупа", "пупа", "loopa", "poopa"})
	if !isThisCommand {
		return false, nil
	}
	p.lock.Lock()
	loopaNews := getLoopaAndPoopaNews("100")
	p.lock.Unlock()

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	if len(loopaNews) <= 0 {
		SendMsg("Семь раз за Пупу и один раз за Лупу")
	} else {
		loopaNews = poopaToKostia(loopaNews)

		_, err := SendMsg(loopaNews)
		if err != nil {
			return true, err
		}
	}
	return true, nil
}

//----------------------------------------------------------------------------------------------------------

// might return an empty string. Depth is int val up to 100
func getLoopaAndPoopaNews(depth string) string {
	parameters := make(map[string]string)
	parameters["count"] = depth // message
	parameters["domain"] = "pupa_and_lupa"

	resp, err := cmn.VkRequest("wall.get", parameters)
	if err != nil {
		log.Printf("wall.get request failed: %s\n", err)
		return ""
	}

	var vkWallResponse cmn.VkWallGetResponse

	err = json.Unmarshal(resp, &vkWallResponse)
	if err != nil {
		log.Printf("Error unmarshalling json response: %s\n", err)
		return ""
	}

	indexes := cmn.InitArrayOfIndexes(len(vkWallResponse.Response.Items))
	if len(indexes) <= 0 {
		log.Println("Error init array of indexes...")
		return ""
	}

	for _, idx := range indexes {
		vkResponseItem := vkWallResponse.Response.Items[idx]

		text := vkResponseItem.Text
		if len(text) <= 0 {
			continue
		}

		// if not found - it might be AD
		if strings.Index(text, "Лупа") == -1 {
			log.Printf("Potential AD detected: %s\n", text)
			continue
		}

		text = strings.Replace(text, "<br>", "\n", -1)

		return text
	}
	return ""
}

//poopaToKostia - превратили Пупу в Костю
func poopaToKostia(msg string) string {
	applyKostia := cmn.Rnd.Intn(12) == 1
	if !applyKostia {
		return msg
	}

	toReplace := map[string]string{
		"Пупа":  "Костя",
		"Пупой": "Костей",
		"Пупу":  "Костю",
		"Пупы":  "Кости",
		"Пупе":  "Косте"}

	for key, val := range toReplace {
		msg = strings.Replace(msg, key, val, -1)
	}

	return msg
}
