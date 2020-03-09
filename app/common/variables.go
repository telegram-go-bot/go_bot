package cmn

//go:generate mockgen -destination=../mocks/mock_randomer.go -package=mocks github.com/telegram-go-bot/go_bot/app/common Randomer

// Randomer - random generater wrapper for test only
type Randomer interface {
	Intn(int) int
}

var (
	// Rnd - pick random number. Interface for tests
	Rnd Randomer
)
