package engine

import (
	"math/rand"
	"time"
)

const (
	gridWidth         = 100
	gridHeight        = 100
	resetToNCreatures = 100
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
	g := make([][]int, gridWidth)
	for i, _ := range g {
		g[i] = make([]int, gridHeight)
	}

	for _, c := range s.Creatures {
		g[c.X][c.Y] = 1
	}

	return g
}

func (s *State) Reset() {
	s.Creatures = make([]*Creature, 0)

	for i := 0; i < resetToNCreatures; i++ {
		c := NewCreature(rand.Intn(gridWidth), rand.Intn(gridHeight))
		s.Creatures = append(s.Creatures, c)
	}
}

func (s *State) Tick() {
	for _, c := range s.Creatures {
		c.Drift()
	}
}
