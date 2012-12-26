package main

import (
	"math/rand"
)

type ReservoirSampler struct {
	reservoir []interface{}
	seen      int
}

func NewSampler(n uint) *ReservoirSampler {
	return &ReservoirSampler{make([]interface{}, n), 0}
}

func (s *ReservoirSampler) Add(x interface{}) {
	n := len(s.reservoir)
	if s.seen < n {
		s.reservoir[s.seen] = x
	} else if rand.Float64() < float64(n)/float64(s.seen) {
		s.reservoir[rand.Intn(int(n))] = x
	}

	s.seen++
}

func (s *ReservoirSampler) Sample() []interface{} {
	return s.reservoir
}
