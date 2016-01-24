package micro

import (
	"bytes"
	"fmt"
)

type CellType byte

const (
	CellTypeNil CellType = iota
	CellTypePlain
	CellTypeAbsorb
	CellTypeAttack
	CellTypeDefend
	CellTypeMove
)

type Cell struct {
	Type  CellType
	Value byte
}

type CellGrid [][]*Cell

func (g CellGrid) String() string {
	var buffer bytes.Buffer
	for _, col := range g {
		for _, c := range col {
			buffer.WriteString(fmt.Sprintf("%v ", c))
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func NewCellGrid(initialWidth, initialHeight int) CellGrid {
	g := make([][]*Cell, initialWidth)
	for i, _ := range g {
		g[i] = make([]*Cell, initialHeight)
	}
	return g
}

func (c Cell) String() string {
	return fmt.Sprintf("cell<%v,%v>", c.Type, c.Value)
}
