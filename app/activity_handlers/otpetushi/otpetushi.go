package otpetushi

import (
	"bytes"
	"image/png"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fogleman/gg"
	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/output"
)

const (
	//TODO: Fix paths somehow
	srcImgPath   = "./app/activity_handlers/otpetushi/resource/src.png"
	srcFontPath  = "./app/activity_handlers/otpetushi/resource/MyriadPro-Regular.ttf"
	srcImgWidth  = 914
	srcImgHeight = 502
	fontSize     = 24
	maxLines     = 5
	maxLineWidth = 360
	lineSpacing  = 1.6
	//text constraints
	topY = 200
)

var (
	img, imgErr = gg.LoadPNG(srcImgPath)
)

type impl struct {
	presenter output.IPresenter
}

// New - constructor
func New(presenter output.IPresenter) impl {
	return impl{
		presenter: presenter}
}

// OnHelp - display help
func (p impl) OnHelp() string {
	return "<b>!отпетуши|petushi</b> <i>replied_to_message</i> - петушим выбраное сообщение"
}

// OnCommand -
func (p impl) OnCommand(item raw.Activity) (bool, error) {

	_, isThisCommand := helpers.IsOnCommand(item.Text, []string{"отпетуши", "петуши", "petushi", "otpetushi"})
	if !isThisCommand || item.RepliedTo == nil {
		return false, nil
	}

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(output.ShowMessageData{ChatID: item.ChatID, Text: message})
	}

	if cmn.Rnd.Intn(100) == 1 {
		_, err := SendMsg(cmn.GetOneMsgFromMany("ti pituh", "ты питух", "ты петух"))
		return true, err
	}

	if imgErr != nil {
		return true, imgErr
	}

	dc, err := initDc()
	if err != nil {
		return true, err
	}

	drawOriginalText(item.RepliedTo.Text, dc)
	drawTranslatedText(messageToKoKo(item.RepliedTo.Text), dc)

	var memImg bytes.Buffer
	png.Encode(&memImg, dc.Image())

	_, err = p.presenter.ShowImage(output.ShowImageData{
		RawImageData:    memImg.Bytes(),
		ShowMessageData: output.ShowMessageData{ChatID: item.ChatID, ReplyToMsgID: item.RepliedTo.MesssageID}})
	return true, nil
}

// have to be called each time new image is created,
// because we write text on DC created
func initDc() (*gg.Context, error) {
	dc := gg.NewContext(srcImgWidth, srcImgHeight)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(srcFontPath, fontSize); err != nil {
		return nil, err
	}
	dc.DrawImage(img, 0, 0)
	return dc, nil
}

func isChar(char rune) bool {
	return !unicode.IsDigit(char) && !unicode.IsPunct(char)
}

func wordToKoKo(word string) string {
	var res strings.Builder

	wordLen := utf8.RuneCountInString(word)

	if wordLen == 0 {
		return ""
	}

	// check if all chars are numbers. Return original string if so
	anyCharExists := strings.IndexFunc(word, isChar)
	if anyCharExists == -1 {
		return word
	}

	if wordLen == 1 {
		return "ко"
	}

	if wordLen == 3 {
		return "коко"
	}

	kokoNum := wordLen / 2
	additionalOes := wordLen % 2

	for idx := 0; idx < kokoNum; idx++ {
		res.WriteString("ко")
	}

	if additionalOes != 0 {
		res.WriteString("о")
	}

	return res.String()
}

func isGoodChar(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}

	if ch >= 'A' && ch <= 'Z' || (ch >= 'a' && ch <= 'z') {
		return true
	}

	if ch >= 'А' && ch <= 'Я' || (ch >= 'а' && ch <= 'я') {
		return true
	}

	if ch == 'Ё' || ch == 'ё' {
		return true
	}

	return false
}

func messageToKoKo(message string) string {

	var res, word strings.Builder

	var flushWord = func() {
		if word.Len() != 0 {
			res.WriteString(wordToKoKo(word.String()))
			word.Reset()
		}
	}

	for _, ch := range message {
		if isGoodChar(ch) {
			word.WriteString(string(ch))
			continue
		}

		flushWord()
		res.WriteString(string(ch))
	}

	flushWord()
	return res.String()
}

func spliToLinesIfNeeded(str string, dc *gg.Context) string {
	lines := dc.WordWrap(str, maxLineWidth)
	if len(lines) > maxLines {
		var maxPossibleString int // cut more than 6 lines to fit "translate window"
		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			maxPossibleString += len(lines[lineIdx])
		}

		str = str[:maxPossibleString]
	}
	return str
}

func drawOriginalText(str string, dc *gg.Context) {
	str = spliToLinesIfNeeded(str, dc)
	dc.DrawStringWrapped(str, 30, topY, 0.0, 0.0, maxLineWidth, lineSpacing, gg.AlignLeft)
}

func drawTranslatedText(str string, dc *gg.Context) {
	str = spliToLinesIfNeeded(str, dc)
	dc.DrawStringWrapped(str, 470, topY, 0.0, 0.0, maxLineWidth, lineSpacing, gg.AlignLeft)
}
