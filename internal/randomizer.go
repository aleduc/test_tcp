package internal

import (
	"math/rand"
	"sync"
	"time"
)

type RandomizerTime struct {
	*rand.Rand
	sync.Mutex
}

// NewRandomizerTime реализует Randomizer на основе текущего времени.
func NewRandomizerTime() *RandomizerTime {
	return &RandomizerTime{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateInt генерация рандомного инта в промежутке [0:n).
func (r *RandomizerTime) GenerateInt(n int) int {
	r.Lock()
	defer r.Unlock()
	return r.Intn(n)
}
