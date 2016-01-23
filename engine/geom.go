package engine

import "galapagia/Godeps/_workspace/src/github.com/dhconnelly/rtreego"

type Positioner interface {
	X() int
	Y() int
}

type Pos struct{ XPos, YPos int }

func (p Pos) X() int { return p.XPos }
func (p Pos) Y() int { return p.YPos }

func positionerToRtreePoint(p Positioner) rtreego.Point {
	return rtreego.Point{float64(p.X()), float64(p.Y())}
}

func nativeCoordsToRtreeRect(x, y, w, h int) *rtreego.Rect {
	r, _ := rtreego.NewRect(rtreego.Point{float64(x), float64(y)}, []float64{float64(w), float64(h)})
	return r
}

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
