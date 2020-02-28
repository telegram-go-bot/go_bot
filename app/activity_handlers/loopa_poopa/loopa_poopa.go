package loopapoopa

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	"github.com/telegram-go-bot/go_bot/app/common/vk"
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
	loopaNews, err := getLoopaAndPoopaNews("100")
	p.lock.Unlock()
	if err != nil {
		return true, err
	}

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	if len(loopaNews) <= 0 {
		SendMsg("Семь раз за Пупу и один раз за Лупу")
		return true, errors.New("len(loopaNews) <= 0")
	}

	loopaNews = poopaToKostia(loopaNews)

	_, err = SendMsg(loopaNews)
	if err != nil {
		return true, err
	}

	return true, nil
}

//----------------------------------------------------------------------------------------------------------

// might return an empty string. Depth is int val up to 100
func getLoopaAndPoopaNews(depth string) (string, error) {
	parameters := make(map[string]string)
	parameters["count"] = depth // message
	parameters["domain"] = "pupa_and_lupa"

	resp, err := vk.Request("wall.get", parameters)
	if err != nil {
		return "", fmt.Errorf("wall.get request failed: %s", err)
	}

	var vkWallResponse vk.WallGetResponse

	err = json.Unmarshal(resp, &vkWallResponse)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling json response: %s", err)
	}

	indexes := vk.InitArrayOfIndexes(len(vkWallResponse.Response.Items))
	if len(indexes) <= 0 {
		return "", errors.New("Error init array of indexes")
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

		return text, nil
	}
	return "", nil
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
