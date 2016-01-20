package engine

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type State struct {
	Creatures []*Creature
}

func NewState() *State {
	cs := make([]*Creature, 0)
	return &State{Creatures: cs}
}

func (s *State) CurrentGrid() [][]int {
	g := make([][]int, 100)
	for i, _ := range g {
		g[i] = make([]int, 100)
	}

	for _, c := range s.Creatures {
		g[c.X][c.Y] = 1
	}

	return g
}

func (s *State) Reset() {
	s.Creatures = make([]*Creature, 0)

	for i := 0; i < 100; i++ {
		c := NewCreature(rand.Intn(100), rand.Intn(100))
		s.Creatures = append(s.Creatures, c)
	}
}
