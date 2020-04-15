package fixlayout

import (
	"errors"
	"sync"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	"github.com/telegram-go-bot/go_bot/app/output"
)

var (
	engToRus = map[string]string{
		"q":  "й",
		"Q":  "Й",
		"w":  "ц",
		"W":  "Ц",
		"e":  "у",
		"E":  "у",
		"r":  "к",
		"R":  "К",
		"t":  "е",
		"T":  "Е",
		"y":  "н",
		"Y":  "Н",
		"u":  "г",
		"U":  "Г",
		"i":  "ш",
		"I":  "Ш",
		"o":  "щ",
		"O":  "Щ",
		"p":  "з",
		"P":  "З",
		"[":  "х",
		"{":  "Х",
		"]":  "ъ",
		"}":  "Ъ",
		"|":  "/",
		"`":  "ё",
		"~":  "Ё",
		"a":  "ф",
		"A":  "Ф",
		"s":  "ы",
		"S":  "Ы",
		"d":  "в",
		"D":  "В",
		"f":  "а",
		"F":  "А",
		"g":  "п",
		"G":  "П",
		"h":  "р",
		"H":  "Р",
		"j":  "о",
		"J":  "О",
		"k":  "л",
		"K":  "Л",
		"l":  "д",
		"L":  "Д",
		";":  "ж",
		":":  "Ж",
		"'":  "э",
		"\"": "Э",
		"z":  "я",
		"Z":  "Я",
		"x":  "ч",
		"X":  "Ч",
		"c":  "с",
		"C":  "С",
		"v":  "м",
		"V":  "М",
		"b":  "и",
		"B":  "И",
		"n":  "т",
		"N":  "Т",
		"m":  "ь",
		"M":  "Ь",
		",":  "б",
		"<":  "Б",
		".":  "ю",
		">":  "Ю",
		"/":  ".",
		"?":  ",",
		"@":  "\"",
		"#":  "№",
		"$":  ";",
		"^":  ":",
		"&":  "?"}

	rusToEng     map[string]string // init it once
	fillRusToEng sync.Once
)

// Impl - implementation
type Impl struct {
	presenter output.IPresenter
}

// New - constructor
func New(presenter output.IPresenter) Impl {

	fillRusToEng.Do(func() {
		rusToEng = make(map[string]string, len(engToRus))
		for key, value := range engToRus {
			rusToEng[value] = key
		}
	})

	return Impl{
		presenter: presenter}
}

// OnHelp - display help
func (p Impl) OnHelp() string {
	return "<b>!fix|fixtext|исправь|фикс</b> <i>replied_to_message</i> - поправим месадж с правильной раскладкой"
}

// OnCommand -
func (p Impl) OnCommand(item raw.Activity) (bool, error) {

	_, isThisCommand := helpers.IsOnCommand(item.Text, []string{"fix", "fixtext", "исправь", "фикс"})
	if !isThisCommand || item.RepliedTo == nil {
		return false, nil
	}

	if len(item.RepliedTo.Text) == 0 {
		return true, nil
	}

	SendMsg := func(message string) (int, error) {
		return p.presenter.ShowMessage(
			output.ShowMessageData{
				ChatID:       item.ChatID,
				Text:         message,
				ReplyToMsgID: item.RepliedTo.MesssageID})
	}

	transformedMsg, err := transformText(item.RepliedTo.Text)
	if err != nil {
		SendMsg(err.Error())
		return true, err
	}

	_, err = SendMsg(transformedMsg)
	if err != nil {
		return true, err
	}

	return true, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isSpecialSymbol(r rune) bool {
	// 0x20..0x40
	if r >= 0x20 && r <= 0x40 {
		return true
	}
	return false
}

// @return (RUS, ENG)
func delectLang(message string) (bool, bool) {
	// lets assume we are on qwerty keyboard

	rus := 0
	eng := 0
	charsEnoughToDecide := min(len(message), 5)

	for itemNo, item := range message {

		if itemNo > charsEnoughToDecide {
			break
		}

		// do not use special symbols in calculations
		if isSpecialSymbol(item) {
			continue
		}

		_, found := rusToEng[string(item)]
		if found {
			rus++
			continue
		}

		_, found = engToRus[string(item)]
		if found {
			eng++
			continue
		}
	}

	// if equal return eng
	if rus == eng && rus != 0 {
		return false, true
	}

	return rus > eng, eng > rus
}

func transformText(message string) (string, error) {
	isRus, isEng := delectLang(message)
	if !isRus && !isEng {
		return "", errors.New("I dont know what language is it: " + message)
	}

	var dict *map[string]string
	if isRus {
		dict = &rusToEng
	} else {
		dict = &engToRus
	}

	var result string
	for _, runeItem := range message {
		char := string(runeItem)

		found, ok := (*dict)[char]
		if !ok {
			result += char
		}
		result += found
	}

	return result, nil
}
