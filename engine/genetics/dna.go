package genetics

import (
	"fmt"
	"galapagia/engine"
	"galapagia/engine/micro"
)

type Sequence []byte

func EncodeBug(s Sequence) *engine.Bug {
	return nil
}

const (
	uniqueMarkerCount = 7
)

const (
	markerNoChild byte = iota
	markerCellTypePlain
	markerCellTypeAbsorb
	markerCellTypeAttack
	markerCellTypeDefend
	markerCellTypeMove
)

type bodyPlanNode struct {
	nodes []*bodyPlanNode // should always have a len == 4
	cell  *micro.Cell
	index int
}

func (bpn bodyPlanNode) String() string {
	return fmt.Sprintf("%v::%v", bpn.cell, bpn.nodes)
}

func debugBodyPlanNode(bpn *bodyPlanNode) string {
	if bpn != nil {
		return fmt.Sprintf(
			"%v@%v::[%v,%v,%v,%v]",
			bpn.cell, bpn.index,
			debugBodyPlanNode(bpn.nodes[0]), debugBodyPlanNode(bpn.nodes[1]), debugBodyPlanNode(bpn.nodes[2]), debugBodyPlanNode(bpn.nodes[3]))
	}
	return "<nil>"
}

func newBodyPlanNode(c *micro.Cell) *bodyPlanNode {
	return &bodyPlanNode{
		nodes: make([]*bodyPlanNode, 4),
		cell:  c,
	}
}

func markerToCellType(m byte) micro.CellType {
	switch m {
	case markerCellTypePlain:
		return micro.CellTypePlain
	case markerCellTypeAbsorb:
		return micro.CellTypeAbsorb
	case markerCellTypeAttack:
		return micro.CellTypeAttack
	case markerCellTypeDefend:
		return micro.CellTypeDefend
	case markerCellTypeMove:
		return micro.CellTypeMove
	default:
		panic("Missing cell type for marker")
	}
}

func sequenceToBodyPlanTree(s Sequence) *bodyPlanNode {
	var root *bodyPlanNode
	stack := newBodyPlanNodeStack()

	for is := 0; is+1 < len(s); is++ {
		m := s[is] % uniqueMarkerCount
		is++

		switch m {
		case markerNoChild:
			if root == nil {
				c := &micro.Cell{Type: micro.CellTypeAbsorb, Value: s[is]}
				root = newBodyPlanNode(c)
				stack.Put(root)
			} else if stack.Empty() {
				return root
			} else {
				current := stack.Peek()
				for current.index == 4 {
					if stack.Empty() {
						return root
					}
					stack.Pop()
					current = stack.Peek()
					if current == nil {
						return root
					}
				}
				current.index++
			}
		default:
			c := &micro.Cell{Type: markerToCellType(m), Value: s[is]}
			if root == nil {
				root = newBodyPlanNode(c)
				stack.Put(root)
			} else if stack.Empty() {
				return root
			} else {
				current := stack.Peek()
				for current.index == 4 {
					if stack.Empty() {
						return root
					}
					stack.Pop()
					current = stack.Peek()
					if current == nil {
						return root
					}
				}
				n := newBodyPlanNode(c)
				current.nodes[current.index] = n
				current.index++
				stack.Put(n)
			}
		}
	}

	return root // TODO should return root
}

func bodyPlanTreeToCellGrid(n *bodyPlanNode) micro.CellGrid {
	// TODO
	g := micro.NewCellGrid(1, 1)
	//cx, cy := 0, 0

	return g
}

func sequenceToCellGrid(s Sequence) micro.CellGrid {
	bpt := sequenceToBodyPlanTree(s)
	cg := bodyPlanTreeToCellGrid(bpt)
	return cg
}
