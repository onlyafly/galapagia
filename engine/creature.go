package engine

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
