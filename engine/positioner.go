package engine

type Positioner interface {
	X() int
	Y() int
}

type Pos struct{ XPos, YPos int }

func (p Pos) X() int { return p.XPos }
func (p Pos) Y() int { return p.YPos }
