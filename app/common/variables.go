package cmn

import (
	"math/rand"
	"time"
)

var (
	// Rnd - init rnd globally
	Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)
