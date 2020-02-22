package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	activityhandlers "github.com/telegram-go-bot/go_bot/app/activity_handlers"
	"github.com/telegram-go-bot/go_bot/app/activity_handlers/goroskop"
	pickfirstorsecond "github.com/telegram-go-bot/go_bot/app/activity_handlers/pick_first_or_second"
	"github.com/telegram-go-bot/go_bot/app/activity_handlers/zagadka"
	collywrapper "github.com/telegram-go-bot/go_bot/app/common/web_scrapper/colly_wrapper"
	"github.com/telegram-go-bot/go_bot/app/domain"
	in "github.com/telegram-go-bot/go_bot/app/input/activities/telegram"
	presenters "github.com/telegram-go-bot/go_bot/app/output/presenters"
	"github.com/telegram-go-bot/go_bot/app/output/views/telegram"
)

var (
	botToken   = os.Getenv("HEROKU_BOT_ID")
	baseURL    = os.Getenv("HEROKU_BASE_URL")
	dbURL      = os.Getenv("DATABASE_URL")
	botUuids   = strings.Split(os.Getenv("BOT_UIDS"), ",")
	vkLogin    = os.Getenv("VK_LOGIN")
	vkPwd      = os.Getenv("VK_PASSWORD")
	botAdminID = os.Getenv("OWNER_ID")
)

var (
	tgView      = telegram.NewTelegramAPIView(botToken)
	tgPresenter = presenters.NewActivityPresenter(tgView)
	tgReader    = in.NewMessageReader(botToken)
	scrapper    = collywrapper.Scrapper{}

	commandHandlers = []activityhandlers.ICommandHandler{
		pickfirstorsecond.New(tgPresenter),
		zagadka.New(tgPresenter, scrapper),
		goroskop.New(tgPresenter, scrapper)}
)

func main() {
	err := domain.InitVKApi(vkLogin, vkPwd)
	if err != nil {
		panic(err)
	}
	tgBot := activityhandlers.NewActivityHandlerFacade(commandHandlers)
	tgBot.ProcessActivities(tgReader)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<3")
}
