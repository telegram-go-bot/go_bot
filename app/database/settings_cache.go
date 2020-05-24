package settings

import raw "github.com/telegram-go-bot/go_bot/app/domain"

var (
	chatInfos = make(map[uint64]*raw.ChatInfo)
	chatUsers = make(map[int]*raw.ChatUser)
)

// ICache - ISettings decorator, cached some results
type ICache interface {
	ISettings
}

// Cache - cached info. High-level Settings Manager.
type Cache struct {
	settings ISettings
}

// NewCache - create new Cache
func NewCache(url string) (ICache, error) {
	cache := &Cache{}
	err := cache.Init(url)
	if err != nil {
		return &Cache{}, err
	}
	return cache, nil
}

// Init - init newely created cache
func (p *Cache) Init(url string) error {
	p.settings = &Settings{}
	return p.settings.Init(url)
}

// AddRecord -
func (p *Cache) AddRecord(newRec interface{}) error {
	return p.settings.AddRecord(newRec)
}

// GetChatInfo -
func (p *Cache) GetChatInfo(chatID int64) (*raw.ChatInfo, error) {
	return p.settings.GetChatInfo(chatID)
}

// GetChatUser -
func (p *Cache) GetChatUser(userID int) (*raw.ChatUser, error) {

	cachedUser, ok := chatUsers[userID]
	if ok {
		return cachedUser, nil
	}

	user, err := p.settings.GetChatUser(userID)
	if err != nil {
		return nil, err
	}

	// cache queried user
	chatUsers[userID] = user
	return user, nil
}
