package engine

import (
	"container/list"
	"fmt"
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
	CreatureList *list.List
	CreatureGrid [][]*Creature
	CellGrid     [][]int // Only for visualization purposes
}

func NewState() *State {
	cl := list.New()

	crg := make([][]*Creature, gridWidth)
	for i, _ := range crg {
		crg[i] = make([]*Creature, gridHeight)
	}

	cellg := make([][]int, gridWidth)
	for i, _ := range cellg {
		cellg[i] = make([]int, gridHeight)
	}

	return &State{
		CreatureList: cl,
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

	for e := s.CreatureList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Creature)
		s.CellGrid[c.xpos][c.ypos] = 1 //TODO
	}
}

func (s *State) LogCreatures() {
	for e := s.CreatureList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Creature)
		fmt.Println("Creature", c)
	}
}

func (s *State) Reset() {
	s.CreatureList = list.New()

	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			s.CreatureGrid[i][j] = nil
		}
	}

	for i := 0; i < resetToNCreatures; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		c := NewCreature(x, y)
		fmt.Println("reset unplaced", c, x, y)
		s.PlaceNewCreature(c, Pos{x, y})
		fmt.Println("reset placed", c, x, y)
	}
}

func closestEmptyPosition(s *State, x, y int) (nx int, ny int, ok bool) {
	if s.CreatureGrid[x][y] == nil {
		return x, y, true
	}

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			nx, ny := restrictToGrid(x+dx, y+dy)
			if s.CreatureGrid[nx][ny] == nil {
				return nx, ny, true
			}
		}
	}
	return -1, -1, false
}

func (s *State) PlaceNewCreature(c *Creature, nearPosition Positioner) (ok bool) {
	s.CreatureList.PushBack(c)

	x, y, ok := closestEmptyPosition(s, nearPosition.X(), nearPosition.Y())
	if !ok {
		return false
	}

	fmt.Println("placing", c, x, y)
	c.xpos = x
	c.ypos = y
	s.CreatureGrid[x][y] = c

	return true
}

func (s *State) Tick() {
	for e := s.CreatureList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Creature)
		s.TickCreature(c, e)
	}
}

func (s *State) TickCreature(c *Creature, celement *list.Element) {
	//TODO s.CheckCreatureVitals(c, celement)
	s.MaybeMoveCreature(c)
	//TODO s.TickCreatureCells(c)
	//TODO s.MaybeReproduceCreature(c)
}

func (s *State) MaybeReproduceCreature(parent *Creature) {
	if parent.Energy > parent.ReproductionCost()*2 {
		child := parent.Reproduce()
		ok := s.PlaceNewCreature(child, parent)
		if !ok {
			// TODO what to do if there is no space for creature?
		}
	}
}

func (s *State) CheckCreatureVitals(c *Creature, celement *list.Element) {
	if c.Energy <= 0 {
		// Kill creature
		s.CreatureGrid[c.xpos][c.ypos] = nil
		s.CreatureList.Remove(celement)
	}
}

func (s *State) TickCreatureCells(c *Creature) {
	// Consumed energy for this tick
	c.Energy -= int(c.CellGrid[0][0].Value / 100)

	// Gained energy for this tick
	c.Energy += int(c.CellGrid[0][0].Value / 10)
}

func (s *State) MaybeMoveCreature(c *Creature) {
	// Should it move?
	if rand.Intn(2) == 0 {
		return // Shouldn't move
	}

	// Where should it move?
	x, y := calcDriftPos(c)
	fmt.Printf("Ticking creature %v to pos %v,%v\n", c, x, y)

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

func restrictToGrid(x, y int) (nx, ny int) {
	nx, ny = x, y
	if x < 0 {
		nx = 0
	}
	if x >= gridWidth {
		nx = gridWidth - 1
	}
	if y < 0 {
		ny = 0
	}
	if y >= gridHeight {
		ny = gridHeight - 1
	}
	return nx, ny
}

func calcDriftPos(p Positioner) (x, y int) {
	dx := rand.Intn(3) - 1 // in range [-1, 1]
	dy := rand.Intn(3) - 1 // in range [-1, 1]
	newX := p.X() + dx
	newY := p.Y() + dy
	return restrictToGrid(newX, newY)
}
