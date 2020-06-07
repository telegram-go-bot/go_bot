package magicquery

import (
	"errors"
	"log"
	"sort"
	"strings"
	"sync"

	helpers "github.com/telegram-go-bot/go_bot/app/activity_handlers/activity_helpers"
	cmn "github.com/telegram-go-bot/go_bot/app/common"
	websearch "github.com/telegram-go-bot/go_bot/app/common/web_search"
	settings "github.com/telegram-go-bot/go_bot/app/database"
	raw "github.com/telegram-go-bot/go_bot/app/domain"
	output "github.com/telegram-go-bot/go_bot/app/output"
)

const (
	valueTypeQuery   = "q"
	valueTypeURL     = "u"
	defaultValueType = valueTypeQuery
	queryDelimiter   = ";;"
)

// MagicQuery - generic magic-> query pair
// Queries are ; separated
type MagicQuery struct {
	ID              uint   `gorm:"AUTO_INCREMENT"`
	Command         string `gorm:"unique_index"` // magic
	DisabledChatIDs string // NOT IMPLEMENTED
	// todo: replace Value+ValueType with json
	Value     string // ' separated Query, or hardcoded url
	ValueType string // string or URL
}

// Values - query values
type Values struct {
	values    []string
	valueType string
	dbID      uint
}

// Impl - implementation
type Impl struct {
	presenter     output.IPresenter
	searcher      websearch.Searcher
	cachedQueries map[string]Values
	cachelocker   *sync.RWMutex
}

// New - constructor
func New(presenter output.IPresenter, searcher websearch.Searcher) Impl {
	var tmp = Impl{presenter: presenter}
	tmp.cachelocker = new(sync.RWMutex)
	tmp.searcher = searcher
	tmp.cachedQueries = make(map[string]Values)
	err := tmp.ReadMagics()
	if err != nil {
		panic("Error reading magics: " + err.Error())
	}
	return tmp
}

func createValuesEntry(value string, valueType string, dbID uint) Values {
	res := Values{valueType: valueType, dbID: dbID}
	res.values = strings.Split(value, queryDelimiter)
	return res
}

// ReadMagics - initial read data from db if any
func (p Impl) ReadMagics() error {
	p.cachelocker.Lock()
	defer p.cachelocker.Unlock()
	var magics []MagicQuery
	log.Printf("ReadMagics::GetHandlerRecords: %v. Settings: %v", magics, settings.Inst())
	err := settings.Inst().GetHandlerRecords(&magics)
	if err != nil {
		return err
	}

	log.Printf("ReadMagics. Magics list: %v", magics)
	for _, savedQuery := range magics {
		p.cachedQueries[savedQuery.Command] =
			createValuesEntry(savedQuery.Value, savedQuery.ValueType, savedQuery.ID)
	}

	return nil
}

// OnHelp - display help
func (p Impl) OnHelp() string {
	return ""
}

func mapTo(magic string, queries Values) MagicQuery {
	res := MagicQuery{}

	var queryStr string
	for _, query := range queries.values {
		if len(queryStr) > 0 {
			queryStr += queryDelimiter
		}
		queryStr += query
	}

	res.ValueType = queries.valueType
	res.Value = queryStr
	res.Command = magic
	res.ID = queries.dbID
	return res
}

func (p Impl) updateRecordInDb(magic string) error {
	p.cachelocker.RLock()
	inMemItem := p.cachedQueries[magic]
	p.cachelocker.RUnlock()

	itemToUpdate := mapTo(magic, inMemItem)

	settings.Inst().UpdateRecord(&itemToUpdate)

	return nil
}

func (p Impl) createRecordInDb(magic string) error {
	p.cachelocker.RLock()
	inMemItem := p.cachedQueries[magic]
	p.cachelocker.RUnlock()

	itemToUpdate := mapTo(magic, inMemItem)

	settings.Inst().AddRecord(&itemToUpdate)

	return nil
}

// addNewMagicQuery - add completely new pair magic->query
func (p Impl) addNewMagicQuery(magic string, values Values) error {
	p.cachelocker.Lock()
	p.cachedQueries[magic] = values
	sort.Strings(p.cachedQueries[magic].values)
	p.cachelocker.Unlock()

	err := p.createRecordInDb(magic)
	if err != nil {
		return err
	}

	return nil
}

// updateMagicQuery - update existing magic with additional query
func (p Impl) updateMagicQuery(magic string, query string) error {
	p.cachelocker.Lock()
	cachedVals, ok := p.cachedQueries[magic]
	if !ok {
		// item have to be present to be updated
		return errors.New("cant update item that doesnt exist")
	}

	cachedVals.values = append(cachedVals.values, query)
	p.cachedQueries[magic] = cachedVals
	p.cachelocker.Unlock()

	p.updateRecordInDb(magic)

	return nil
}

type parsedChatParameter struct {
	magic string // command
	typ   string // type
	val   string // query
}

// param format - <magic> <val> <type> (optional)
func parseAdd(text string) (*parsedChatParameter, error) {
	res := strings.Split(text, queryDelimiter)
	if len(res) == 3 {
		return &parsedChatParameter{
			magic: strings.ToLower(res[0]),
			val:   strings.ToLower(res[1]),
			typ:   res[2]}, nil
	} else if len(res) == 2 {
		return &parsedChatParameter{
			magic: strings.ToLower(res[0]),
			val:   strings.ToLower(res[1]),
			typ:   defaultValueType}, nil
	}

	return nil, errors.New("Invalid param")
}

func (p Impl) onAddMagic(param string) error {

	parsed, err := parseAdd(param)
	if err != nil {
		return err
	}

	values := createValuesEntry(parsed.val, parsed.typ, 0)
	return p.addNewMagicQuery(parsed.magic, values)
}

// param format - <magic> <val>
func (p Impl) onUpdateMagic(param string) error {

	parsed, err := parseAdd(param)
	if err != nil {
		return err
	}

	return p.updateMagicQuery(parsed.magic, parsed.val)
}

func (p Impl) isOnDynamicCommand(text string) (*Values, bool) {
	if len(text) == 0 {
		return nil, false
	}

	textToCheck := strings.ToLower(text)

	p.cachelocker.RLock()
	defer p.cachelocker.RUnlock()

	for key, val := range p.cachedQueries {
		cmdSymbIdx := strings.Index(textToCheck, key)
		if cmdSymbIdx != -1 {
			return &val, true
		}
	}

	return nil, false
}

// OnCommand -
func (p Impl) OnCommand(item raw.Activity) (bool, error) {
	command, isThisCommand := helpers.IsOnCommand(item.Text, []string{"add_magic"})
	if isThisCommand {
		// here we should process internal add_magic\update_magic etc. commands
		err := p.onAddMagic(command)
		return true, err
	}

	command, isThisCommand = helpers.IsOnCommand(item.Text, []string{"update_magic"})
	if isThisCommand {
		// here we should process internal add_magic\update_magic etc. commands
		err := p.onUpdateMagic(command)
		return true, err
	}

	values, isThisCommand := p.isOnDynamicCommand(item.Text)
	if !isThisCommand {
		return true, nil
	}

	pickN := cmn.Rnd.Intn(len(values.values))
	query := values.values[pickN]

	var resultURL string
	if values.valueType == valueTypeURL {
		resultURL = query
	} else {
		// else search for image url
		images := p.searcher.SearchImage(query, 10)
		if len(images) == 0 {
			return false, nil
		}

		pickN := cmn.Rnd.Intn(len(images))
		resultURL = images[pickN]
	}

	p.presenter.ShowImage(output.ShowImageData{
		ImageURL:        resultURL,
		ShowMessageData: output.ShowMessageData{ChatID: item.ChatID}})

	return true, nil
}
