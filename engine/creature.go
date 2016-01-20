package engine

import "fmt"

type CellType byte

const (
	cellTypePlain CellType = iota
	cellTypeMovement
	cellTypeDefence
	cellTypeAttack
	cellTypeAbsorb
)

type Cell struct {
	Type  CellType
	Value byte
}

type Creature struct {
	CellGrid [][]Cell
	Width    int
	Height   int
	xpos     int
	ypos     int
	Energy   int
}

func NewCreature(x, y int) *Creature {
	w := 1
	h := 1
	g := make([][]Cell, w)
	for i, _ := range g {
		g[i] = make([]Cell, h)
		g[i][0] = Cell{cellTypeAbsorb, 127}
	}
	return &Creature{
		CellGrid: g,
		Width:    w,
		Height:   h,
		xpos:     x,
		ypos:     y,
		Energy:   100, // TODO
	}
}

func (c *Creature) X() int {
	return c.xpos
}

func (c *Creature) Y() int {
	return c.ypos
}

func (c *Creature) String() string {
	return fmt.Sprintf("c<e=%v>", c.Energy)
}
