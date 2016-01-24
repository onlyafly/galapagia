package micro

import "fmt"

type CellType byte

const (
	CellTypePlain CellType = iota
	CellTypeAbsorb
	CellTypeAttack
	CellTypeDefend
	CellTypeMove
)

type Cell struct {
	Type  CellType
	Value byte
}

type CellGrid [][]Cell

func (c CellGrid) String() string {
	// TODO
	return fmt.Sprintf("1")
}

func NewCellGrid(initialWidth, initialHeight int) CellGrid {
	g := make([][]Cell, initialWidth)
	for i, _ := range g {
		g[i] = make([]Cell, initialHeight)
	}
	return g
}

func (c Cell) String() string {
	return fmt.Sprintf("cell<%v,%v>", c.Type, c.Value)
}
