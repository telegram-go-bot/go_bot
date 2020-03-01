package cmn

// Randomer - random generater wrapper for test only
type Randomer interface {
	Intn(int) int
}

var (
	// Rnd - pick random number. Interface for tests
	Rnd Randomer
)
