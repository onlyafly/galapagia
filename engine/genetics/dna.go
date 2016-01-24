package genetics

import (
	"fmt"
	"galapagia/engine/micro"
	"math/rand"
)

type Sequence []byte

func RandomSequence(pairs int) Sequence {
	bs := make([]byte, pairs*2)
	for i, _ := range bs {
		bs[i] = byte(rand.Intn(255))
	}
	return bs
}

const (
	uniqueMarkerCount = 6
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
	depth int
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

func sequenceToBodyPlanTree(s Sequence) (root *bodyPlanNode, deepestNodeDepth int) {
	stack := newBodyPlanNodeStack()
	deepestNodeDepth = 1

	for is := 0; is+1 < len(s); is++ {
		m := s[is] % uniqueMarkerCount
		is++

		switch m {
		case markerNoChild:
			if root == nil {
				c := &micro.Cell{Type: micro.CellTypeAbsorb, Value: s[is]}
				root = newBodyPlanNode(c)
				root.depth = 1
				stack.Put(root)
			} else if stack.Empty() {
				return
			} else {
				current := stack.Peek()
				for current.index == 4 {
					if stack.Empty() {
						return
					}
					stack.Pop()
					current = stack.Peek()
					if current == nil {
						return
					}
				}
				current.index++
			}
		default:
			c := &micro.Cell{Type: markerToCellType(m), Value: s[is]}
			if root == nil {
				root = newBodyPlanNode(c)
				root.depth = 1
				stack.Put(root)
			} else if stack.Empty() {
				return
			} else {
				current := stack.Peek()
				for current.index == 4 {
					if stack.Empty() {
						return
					}
					stack.Pop()
					current = stack.Peek()
					if current == nil {
						return
					}
				}
				n := newBodyPlanNode(c)
				n.depth = current.depth + 1
				if n.depth > deepestNodeDepth {
					deepestNodeDepth = n.depth
				}
				current.nodes[current.index] = n
				current.index++
				stack.Put(n)
			}
		}
	}

	return
}

func mapBodyPlanNodeToGrid(n *bodyPlanNode, g *micro.CellGrid, cx, cy int) {
	if n == nil {
		return
	}
	if (*g)[cx][cy] == nil {
		(*g)[cx][cy] = n.cell
	}
	mapBodyPlanNodeToGrid(n.nodes[0], g, cx, cy-1)
	mapBodyPlanNodeToGrid(n.nodes[1], g, cx+1, cy)
	mapBodyPlanNodeToGrid(n.nodes[2], g, cx, cy+1)
	mapBodyPlanNodeToGrid(n.nodes[3], g, cx-1, cy)
}

func bodyPlanTreeToCellGrid(n *bodyPlanNode, depth int) micro.CellGrid {
	// TODO
	gridSize := 2*depth - 1
	centerPos := gridSize / 2
	g := micro.NewCellGrid(gridSize, gridSize)
	mapBodyPlanNodeToGrid(n, &g, centerPos, centerPos)
	return g
}

func SequenceToCellGrid(s Sequence) micro.CellGrid {
	bpt, depth := sequenceToBodyPlanTree(s)
	cg := bodyPlanTreeToCellGrid(bpt, depth)
	return cg
}
