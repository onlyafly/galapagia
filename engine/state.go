package engine

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"

	"galapagia/engine/micro"

	"galapagia/Godeps/_workspace/src/github.com/dhconnelly/rtreego"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type State struct {
	BugList      *list.List
	BugGrid      [][]*Bug
	BugTree      *rtreego.Rtree
	GridWidth    int
	GridHeight   int
	ResetToNBugs int
}

func NewState(gridWidth, gridHeight int) *State {
	bl := list.New()

	bg := make([][]*Bug, gridWidth)
	for i, _ := range bg {
		bg[i] = make([]*Bug, gridHeight)
	}

	bt := rtreego.NewTree(2, 25, 50)

	return &State{
		BugList:    bl,
		BugGrid:    bg,
		BugTree:    bt,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}
}

func (s *State) Reset(initialBugCount int) {
	s.BugList = list.New()

	for i := 0; i < s.GridWidth; i++ {
		for j := 0; j < s.GridHeight; j++ {
			s.BugGrid[i][j] = nil
		}
	}

	s.BugTree = rtreego.NewTree(2, 25, 50)

	for i := 0; i < initialBugCount; i++ {
		x := rand.Intn(s.GridWidth)
		y := rand.Intn(s.GridHeight)
		b := NewRandomBug(x, y)
		s.PlaceNewBug(b, Pos{x, y})
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
		fmt.Printf("b,b.Genome,b.Body %v %v %v\n", b, b.Genome, b.Body)

		for i := 0; i < b.W(); i++ {
			for j := 0; j < b.H(); j++ {
				cell := b.Body[i][j]
				if cell != nil {
					cellg[b.xpos+i][b.ypos+j] = int(cell.Type)
				}
			}
		}
	}

	return cellg
}

func (s *State) LogBugs() {
	for e := s.BugList.Front(); e != nil; e = e.Next() {
		b := e.Value.(*Bug)
		fmt.Println("Bug", b)
	}
}

// TODO eventually update this to check for actual intersections of cells, not just the bounding box
func hasIntersections(rt *rtreego.Rtree, r *rtreego.Rect) bool {
	results := rt.SearchIntersect(r)
	return len(results) > 0
}

// Check if there any entities in rt that intersect the Rect r, ignoring the spatial entity s.
func hasIntersectionsOtherThanSelf(rt *rtreego.Rtree, r *rtreego.Rect, self rtreego.Spatial) bool {
	results := rt.SearchIntersect(r)
	switch {
	case len(results) == 0:
		return false // intersects nothing
	case len(results) == 1:
		if results[0] == self {
			return false // intersects only itself
		} else {
			return true // intersects another entity
		}
	default:
		return true // intersects more than one entity
	}
}

func closestEmptyPosition(s *State, r *rtreego.Rect) (nx int, ny int, ok bool) {
	x := int(r.PointCoord(0))
	y := int(r.PointCoord(1))
	w := int(r.LengthsCoord(0))
	h := int(r.LengthsCoord(1))

	// These loops try to place the creature at distances at least far enough away to
	// avoid intersecting it
	for dx := -w; dx <= w; dx += w {
		for dy := -h; dy <= h; dy += h {
			nx, ny := restrictToGrid(s, x+dx, y+dy, w, h)

			restrictedRect, _ := rtreego.NewRect(rtreego.Point{float64(nx), float64(ny)}, []float64{r.LengthsCoord(0), float64(r.LengthsCoord(1))})

			if !hasIntersections(s.BugTree, restrictedRect) {
				return nx, ny, true
			}
		}
	}
	return -1, -1, false
}

func (s *State) PlaceNewBug(b *Bug, nearPosition Positioner) (ok bool) {
	desiredRect, _ := rtreego.NewRect(positionerToRtreePoint(nearPosition), []float64{float64(b.W()), float64(b.H())})
	x, y, ok := closestEmptyPosition(s, desiredRect)
	if !ok {
		return false
	}

	s.BugList.PushBack(b)
	b.xpos = x
	b.ypos = y
	s.BugGrid[x][y] = b
	s.BugTree.Insert(b)

	return true
}

func (s *State) Tick() {
	for e := s.BugList.Front(); e != nil; e = e.Next() {
		b := e.Value.(*Bug)
		s.TickBug(b, e)
	}
}

func (s *State) TickBug(b *Bug, celement *list.Element) {
	s.CheckBugVitals(b, celement)
	s.MaybeMoveBug(b)
	s.TickBugCells(b)
	s.MaybeReproduceBug(b)
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

func (s *State) CheckBugVitals(b *Bug, celement *list.Element) {
	if b.Energy <= 0 {
		// Kill bug
		s.BugGrid[b.xpos][b.ypos] = nil
		s.BugList.Remove(celement)
	}
}

func (s *State) TickBugCells(b *Bug) {
	for _, column := range b.Body {
		for _, cell := range column {
			if cell != nil {
				// Consumed energy for this tick
				b.Energy -= int(cell.Value / 100)

				switch cell.Type {
				case micro.CellTypeAbsorb:
					// Gained energy for this tick
					b.Energy += int(cell.Value / 10)
				default:
					// Nothing
				}
			}
		}
	}
}

func (s *State) MaybeMoveBug(b *Bug) {
	// Should it move?
	if rand.Intn(2) == 0 {
		return // Shouldn't move
	}

	// Where should it try to move?
	x, y := calcRandomDriftPos(s, b)

	// Can it move there?
	potentialDriftRect := nativeCoordsToRtreeRect(x, y, b.width, b.height)
	if hasIntersectionsOtherThanSelf(s.BugTree, potentialDriftRect, b) {
		return // Can't move
	}

	// Remove the bug from its old location
	s.BugGrid[b.xpos][b.ypos] = nil // TODO: is BugGrid needed at all???
	s.BugTree.Delete(b)

	// Update the bug's location
	b.xpos = x
	b.ypos = y

	// Insert the bug at its new location
	s.BugGrid[x][y] = b
	s.BugTree.Insert(b)
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

func calcRandomDriftPos(gridSize Sizer, s SquareOnPlane) (x, y int) {
	dx := rand.Intn(3) - 1 // in range [-1, 1]
	dy := rand.Intn(3) - 1 // in range [-1, 1]
	newX := s.X() + dx
	newY := s.Y() + dy
	return restrictToGrid(gridSize, newX, newY, s.W(), s.H())
}
