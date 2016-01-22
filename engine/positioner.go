package engine

type Positioner interface {
	X() int
	Y() int
}

type Pos struct{ XPos, YPos int }

func (p Pos) X() int { return p.XPos }
func (p Pos) Y() int { return p.YPos }

type SquareOnPlane interface {
	Positioner
	W() int
	H() int
}

type Sizer interface {
	W() int
	H() int
}

type Size struct{ w, h int }

func (s Size) W() int { return s.w }
func (s Size) H() int { return s.h }
