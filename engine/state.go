package engine

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type State struct {
	BugList      *list.List
	BugGrid      [][]*Bug
	GridWidth    int
	GridHeight   int
	ResetToNBugs int
}

func NewState(gridWidth, gridHeight int) *State {
	cl := list.New()

	crg := make([][]*Bug, gridWidth)
	for i, _ := range crg {
		crg[i] = make([]*Bug, gridHeight)
	}

	return &State{
		BugList:    cl,
		BugGrid:    crg,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}
}

func (s *State) W() int { return s.GridWidth }
func (s *State) H() int { return s.GridHeight }

func (s *State) CurrentCellGrid() [][]int {
	cellg := make([][]int, s.GridWidth)
	for i, _ := range cellg {
		cellg[i] = make([]int, s.GridHeight)
	}

	for e := s.BugList.Front(); e != nil; e = e.Next() {
		b := e.Value.(*Bug)

		for i := 0; i < b.W(); i++ {
			for j := 0; j < b.H(); j++ {
				cellg[b.xpos+i][b.ypos+j] = int(b.CellGrid[i][j].Type)
			}
		}
	}

	return cellg
}

func (s *State) LogBugs() {
	for e := s.BugList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Bug)
		fmt.Println("Bug", c)
	}
}

func (s *State) Reset(initialBugCount int) {
	s.BugList = list.New()

	for i := 0; i < s.GridWidth; i++ {
		for j := 0; j < s.GridHeight; j++ {
			s.BugGrid[i][j] = nil
		}
	}

	for i := 0; i < initialBugCount; i++ {
		x := rand.Intn(s.GridWidth)
		y := rand.Intn(s.GridHeight)
		c := NewBug(x, y)
		s.PlaceNewBug(c, Pos{x, y})
	}
}

func closestEmptyPosition(s *State, x, y, w, h int) (nx int, ny int, ok bool) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			nx, ny := restrictToGrid(s, x+dx, y+dy, w, h)
			// TODO: this needs to check intersection with other bugs
			if s.BugGrid[nx][ny] == nil {
				return nx, ny, true
			}
		}
	}
	return -1, -1, false
}

func (s *State) PlaceNewBug(c *Bug, nearPosition Positioner) (ok bool) {
	x, y, ok := closestEmptyPosition(s, nearPosition.X(), nearPosition.Y(), c.W(), c.H())
	if !ok {
		return false
	}

	if x > 97 || y > 97 {
		fmt.Println("PlaceNewBug UH OH", c)
	}

	s.BugList.PushBack(c)
	c.xpos = x
	c.ypos = y
	s.BugGrid[x][y] = c

	return true
}

func (s *State) Tick() {
	for e := s.BugList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Bug)
		s.TickBug(c, e)
	}
}

func (s *State) TickBug(c *Bug, celement *list.Element) {
	s.CheckBugVitals(c, celement)
	s.MaybeMoveBug(c)
	s.TickBugCells(c)
	s.MaybeReproduceBug(c)
}

func (s *State) MaybeReproduceBug(parent *Bug) {
	if parent.Energy > parent.ReproductionCost()*2 {
		child := parent.Reproduce()
		ok := s.PlaceNewBug(child, parent)
		if !ok {
			// TODO what to do if there is no space for bug?
		}
	}
}

func (s *State) CheckBugVitals(c *Bug, celement *list.Element) {
	if c.Energy <= 0 {
		// Kill bug
		s.BugGrid[c.xpos][c.ypos] = nil
		s.BugList.Remove(celement)
	}
}

func (s *State) TickBugCells(c *Bug) {
	// TODO tick the cells in the actual bug

	// Consumed energy for this tick
	c.Energy -= int(c.CellGrid[0][0].Value / 100)

	// Gained energy for this tick
	c.Energy += int(c.CellGrid[0][0].Value / 10)
}

func (s *State) MaybeMoveBug(c *Bug) {
	// Should it move?
	if rand.Intn(2) == 0 {
		return // Shouldn't move
	}

	// Where should it move?
	x, y := calcDriftPos(s, c)

	// Can it move there?
	if s.BugGrid[x][y] != nil {
		return // Can't move
	}

	// Move it there

	// ORDERING: must update the bug's position after removing the bug from the grid
	s.BugGrid[c.xpos][c.ypos] = nil
	c.xpos = x
	c.ypos = y
	s.BugGrid[x][y] = c
}

func restrictToGrid(gridSize Sizer, x, y, w, h int) (nx, ny int) {
	nx, ny = x, y
	if x < 0 {
		nx = 0
	}
	if x+w > gridSize.W() {
		nx = gridSize.W() - w
	}
	if y < 0 {
		ny = 0
	}
	if y+h > gridSize.H() {
		ny = gridSize.H() - h
	}
	return nx, ny
}

func calcDriftPos(gridSize Sizer, s SquareOnPlane) (x, y int) {
	dx := rand.Intn(3) - 1 // in range [-1, 1]
	dy := rand.Intn(3) - 1 // in range [-1, 1]
	newX := s.X() + dx
	newY := s.Y() + dy
	return restrictToGrid(gridSize, newX, newY, s.W(), s.H())
}
