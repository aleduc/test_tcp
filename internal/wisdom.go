package internal

type Randomizer interface {
	GenerateInt(n int) int
}

type WisdomHT struct {
	randomizer Randomizer
	cnt        int
	wisdoms    []string
}

func NewWisdomHT(randomizer Randomizer) *WisdomHT {
	return &WisdomHT{randomizer: randomizer, cnt: 3, wisdoms: []string{"wisdom1", "wisdom2", "wisdom3"}}
}

func (w *WisdomHT) Get() (string, error) {
	return w.wisdoms[w.randomizer.GenerateInt(w.cnt)], nil
}
