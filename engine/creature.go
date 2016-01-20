package engine

import "math/rand"

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
	X        int
	Y        int
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
		X:        x,
		Y:        y,
	}
}

func (c *Creature) Drift() {
	dx := rand.Intn(3) - 1 // in range [-1, 1]
	dy := rand.Intn(3) - 1 // in range [-1, 1]
	c.X += dx
	c.Y += dy

	if c.X < 0 {
		c.X = 0
	}
	if c.X >= gridWidth {
		c.X = gridWidth - 1
	}
	if c.Y < 0 {
		c.Y = 0
	}
	if c.Y >= gridHeight {
		c.Y = gridHeight - 1
	}
}
