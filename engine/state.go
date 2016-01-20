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
	Creatures    []*Creature
	CreatureGrid [][]*Creature
	CellGrid     [][]int // Only for visualization purposes
}

func NewState() *State {
	cs := make([]*Creature, 0)

	crg := make([][]*Creature, gridWidth)
	for i, _ := range crg {
		crg[i] = make([]*Creature, gridHeight)
	}

	cellg := make([][]int, gridWidth)
	for i, _ := range cellg {
		cellg[i] = make([]int, gridHeight)
	}

	return &State{
		Creatures:    cs,
		CreatureGrid: crg,
		CellGrid:     cellg,
	}
}

func (s *State) CurrentCellGrid() [][]int {
	s.UpdateCellGrid()
	return s.CellGrid
}

func (s *State) UpdateCellGrid() {
	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			s.CellGrid[i][j] = 0
		}
	}

	for _, c := range s.Creatures {
		s.CellGrid[c.X][c.Y] = 1 //TODO
	}
}

func (s *State) Reset() {
	s.Creatures = make([]*Creature, 0)

	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			s.CreatureGrid[i][j] = nil
		}
	}

	for i := 0; i < resetToNCreatures; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		c := NewCreature(x, y)

		s.Creatures = append(s.Creatures, c)

		s.CreatureGrid[x][y] = c
	}
}

func (s *State) Tick() {
	for _, c := range s.Creatures {
		c.Drift()
	}
}
