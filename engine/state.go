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
		s.CellGrid[c.xpos][c.ypos] = 1 //TODO
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
		s.TickCreature(c)
	}
}

func (s *State) TickCreature(c *Creature) {
	// Should it move?
	if rand.Intn(2) == 0 {
		return // Shouldn't move
	}

	// Where should it move?
	x, y := calcDriftPos(c)

	// Can it move there?
	if s.CreatureGrid[x][y] != nil {
		return // Can't move
	}

	// Move it there

	// ORDERING: must update the creature's position after removing the creature from the grid
	s.CreatureGrid[c.xpos][c.ypos] = nil
	c.xpos = x
	c.ypos = y
	s.CreatureGrid[x][y] = c
}

func calcDriftPos(p Positioner) (x, y int) {
	dx := rand.Intn(3) - 1 // in range [-1, 1]
	dy := rand.Intn(3) - 1 // in range [-1, 1]
	newX := p.X() + dx
	newY := p.Y() + dy

	if newX < 0 {
		newX = 0
	}
	if newX >= gridWidth {
		newX = gridWidth - 1
	}
	if newY < 0 {
		newY = 0
	}
	if newY >= gridHeight {
		newY = gridHeight - 1
	}

	return newX, newY
}
