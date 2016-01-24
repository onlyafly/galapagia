package engine

import (
	"fmt"

	"galapagia/engine/micro"

	"galapagia/Godeps/_workspace/src/github.com/dhconnelly/rtreego"
)

type Bug struct {
	Body   micro.CellGrid
	width  int
	height int
	xpos   int
	ypos   int
	Energy int
}

func NewBug(x, y int) *Bug {
	w := 3
	h := 3

	g := micro.NewCellGrid(w, h)

	g[1][0] = micro.Cell{micro.CellTypeAttack, 127}
	g[0][1] = micro.Cell{micro.CellTypeAttack, 127}
	g[1][1] = micro.Cell{micro.CellTypeAbsorb, 127}
	g[2][1] = micro.Cell{micro.CellTypeAttack, 127}
	g[1][2] = micro.Cell{micro.CellTypeAttack, 127}

	return &Bug{
		Body:   g,
		width:  w,
		height: h,
		xpos:   x,
		ypos:   y,
		Energy: 100, // TODO
	}
}

func (c *Bug) X() int {
	return c.xpos
}

func (c *Bug) Y() int {
	return c.ypos
}

func (c *Bug) W() int {
	return c.width
}

func (c *Bug) H() int {
	return c.height
}

func (c *Bug) String() string {
	return fmt.Sprintf("c<e=%v,x=%v,y=%v>", c.Energy, c.xpos, c.ypos)
}

func (c *Bug) ReproductionCost() int {
	// TODO calculate this based on the cells in the bug
	return 1000
}

// Required for the Spatial interface from rtreego
func (b *Bug) Bounds() *rtreego.Rect {
	r, _ := rtreego.NewRect(
		rtreego.Point{float64(b.xpos), float64(b.ypos)},
		[]float64{float64(b.width), float64(b.height)},
	)
	return r
}

func (parent *Bug) Reproduce() *Bug {
	parent.Energy -= parent.ReproductionCost()

	child := NewBug(parent.xpos, parent.ypos)
	child.Energy = parent.Energy / 2
	parent.Energy = parent.Energy / 2

	return child
}
