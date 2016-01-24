package genetics

import (
	"galapagia/engine/micro"
	"testing"
)

func assertEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("ACTUAL != EXPECTED\nACTUAL:   <%v>\nEXPECTED: <%v>\n", actual, expected)
	}
}

func assertEqualInt(t *testing.T, actual, expected int) {
	if actual != expected {
		t.Errorf("ACTUAL != EXPECTED\nACTUAL:   %q\nEXPECTED: %q\n", actual, expected)
	}
}

func Test_bodyPlanTreeToCellGrid(t *testing.T) {
	var s []byte
	var g micro.CellGrid

	s = []byte{markerNoChild, 200}
	g = bodyPlanTreeToCellGrid(sequenceToBodyPlanTree(s))
	assertEqual(t, g.String(), "cell<1,200> \n")

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerNoChild, 2,
		markerNoChild, 3,
		markerNoChild, 4,
		markerNoChild, 5,
	}
	g = bodyPlanTreeToCellGrid(sequenceToBodyPlanTree(s))
	assertEqual(t, g.String(), "cell<1,1> \n")

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerNoChild, 2,
		markerNoChild, 3,
		markerCellTypeAbsorb, 4,
		markerNoChild, 41,
		markerNoChild, 42,
		markerNoChild, 43,
		markerNoChild, 44,
		markerNoChild, 5,
	}
	g = bodyPlanTreeToCellGrid(sequenceToBodyPlanTree(s))
	assertEqual(t, g.String(), "<nil> <nil> <nil> \n<nil> cell<1,1> cell<1,4> \n<nil> <nil> <nil> \n")

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerCellTypeAbsorb, 2,
		markerCellTypeAbsorb, 21,
		markerNoChild, 21,
		markerNoChild, 21,
		markerNoChild, 21,
		markerNoChild, 21,
		markerNoChild, 22,
		markerNoChild, 23,
		markerNoChild, 24,
		markerCellTypeAbsorb, 3,
		markerNoChild, 31,
		markerCellTypeAbsorb, 32,
		markerNoChild, 32,
		markerNoChild, 32,
		markerNoChild, 32,
		markerNoChild, 32,
		markerNoChild, 33,
		markerNoChild, 34,
		markerCellTypeAbsorb, 4,
		markerNoChild, 41,
		markerNoChild, 42,
		markerCellTypeAbsorb, 43,
		markerNoChild, 43,
		markerNoChild, 43,
		markerNoChild, 43,
		markerNoChild, 43,
		markerNoChild, 44,
		markerCellTypeAbsorb, 5,
		markerNoChild, 51,
		markerNoChild, 52,
		markerNoChild, 53,
		markerCellTypeAbsorb, 54,
		markerNoChild, 54,
		markerNoChild, 54,
		markerNoChild, 54,
		markerNoChild, 54,
	}
	g = bodyPlanTreeToCellGrid(sequenceToBodyPlanTree(s))
	assertEqual(t, g.String(), "<nil> <nil> <nil> \n<nil> cell<1,1> cell<1,4> \n<nil> <nil> <nil> \n")
}

func Test_sequenceToBodyPlanTree(t *testing.T) {
	var s []byte
	var d int // depth
	var bpt *bodyPlanNode

	s = []byte{markerNoChild, 200}
	bpt, d = sequenceToBodyPlanTree(s)
	assertEqual(t, bpt.String(), "cell<1,200>::[<nil> <nil> <nil> <nil>]")
	assertEqualInt(t, d, 1)

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerNoChild, 2,
		markerNoChild, 3,
		markerNoChild, 4,
		markerNoChild, 5,
	}
	bpt, d = sequenceToBodyPlanTree(s)
	assertEqual(t, bpt.String(), "cell<1,1>::[<nil> <nil> <nil> <nil>]")
	assertEqualInt(t, d, 1)

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerNoChild, 2,
		markerNoChild, 3,
		markerCellTypeAbsorb, 4,
		markerNoChild, 41,
		markerNoChild, 42,
		markerNoChild, 43,
		markerNoChild, 44,
		markerNoChild, 5,
	}
	bpt, d = sequenceToBodyPlanTree(s)
	assertEqual(t, bpt.String(), "cell<1,1>::[<nil> <nil> cell<1,4>::[<nil> <nil> <nil> <nil>] <nil>]")
	assertEqualInt(t, d, 2)

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerNoChild, 2,
		markerNoChild, 3,
		markerCellTypeAbsorb, 4,
		markerNoChild, 41,
		markerCellTypeAbsorb, 42,
		markerNoChild, 21,
		markerNoChild, 22,
		markerNoChild, 23,
		markerNoChild, 24,
		markerNoChild, 43,
		markerNoChild, 44,
		markerNoChild, 5,
	}
	bpt, d = sequenceToBodyPlanTree(s)
	assertEqual(t, bpt.String(), "cell<1,1>::[<nil> <nil> cell<1,4>::[<nil> cell<1,42>::[<nil> <nil> <nil> <nil>] <nil> <nil>] <nil>]")
	assertEqualInt(t, d, 3)

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerCellTypeAbsorb, 2,
		markerNoChild, 21,
		markerNoChild, 22,
		markerNoChild, 23,
		markerNoChild, 24,
		markerCellTypeAbsorb, 3,
		markerNoChild, 31,
		markerNoChild, 32,
		markerNoChild, 33,
		markerNoChild, 34,
		markerCellTypeAbsorb, 4,
		markerNoChild, 41,
		markerNoChild, 42,
		markerNoChild, 43,
		markerNoChild, 44,
		markerCellTypeAbsorb, 5,
		markerNoChild, 51,
		markerNoChild, 52,
		markerNoChild, 53,
		markerNoChild, 54,
	}
	bpt, d = sequenceToBodyPlanTree(s)
	assertEqual(t, bpt.String(), "cell<1,1>::[cell<1,2>::[<nil> <nil> <nil> <nil>] cell<1,3>::[<nil> <nil> <nil> <nil>] cell<1,4>::[<nil> <nil> <nil> <nil>] cell<1,5>::[<nil> <nil> <nil> <nil>]]")
	assertEqualInt(t, d, 2)

	s = []byte{
		markerCellTypeAbsorb, 1,
		markerCellTypeAbsorb, 2,
		markerCellTypeAbsorb, 21,
		markerNoChild, 21,
		markerNoChild, 21,
		markerNoChild, 21,
		markerNoChild, 21,
		markerNoChild, 22,
		markerNoChild, 23,
		markerNoChild, 24,
		markerCellTypeAbsorb, 3,
		markerNoChild, 31,
		markerCellTypeAbsorb, 32,
		markerNoChild, 32,
		markerNoChild, 32,
		markerNoChild, 32,
		markerNoChild, 32,
		markerNoChild, 33,
		markerNoChild, 34,
		markerCellTypeAbsorb, 4,
		markerNoChild, 41,
		markerNoChild, 42,
		markerCellTypeAbsorb, 43,
		markerNoChild, 43,
		markerNoChild, 43,
		markerNoChild, 43,
		markerNoChild, 43,
		markerNoChild, 44,
		markerCellTypeAbsorb, 5,
		markerNoChild, 51,
		markerNoChild, 52,
		markerNoChild, 53,
		markerCellTypeAbsorb, 54,
		markerNoChild, 54,
		markerNoChild, 54,
		markerNoChild, 54,
		markerNoChild, 54,
	}
	bpt, d = sequenceToBodyPlanTree(s)
	assertEqual(t, bpt.String(), "cell<1,1>::[cell<1,2>::[cell<1,21>::[<nil> <nil> <nil> <nil>] <nil> <nil> <nil>] cell<1,3>::[<nil> cell<1,32>::[<nil> <nil> <nil> <nil>] <nil> <nil>] cell<1,4>::[<nil> <nil> cell<1,43>::[<nil> <nil> <nil> <nil>] <nil>] cell<1,5>::[<nil> <nil> <nil> cell<1,54>::[<nil> <nil> <nil> <nil>]]]")
	assertEqualInt(t, d, 3)
}
