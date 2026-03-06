package random

import "math/rand"

type Generator interface {
	Intn(n int) int
}

type defaultGenerator struct{}

func (defaultGenerator) Intn(n int) int {
	return rand.Intn(n)
}

func NewGenerator(seed int64) Generator {
	if seed < 0 {
		return defaultGenerator{}
	}
	return rand.New(rand.NewSource(seed))
}
