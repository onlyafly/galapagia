package genetics

type bodyPlanNodeStack struct{ v []*bodyPlanNode }

func newBodyPlanNodeStack() *bodyPlanNodeStack {
	return &bodyPlanNodeStack{
		v: make([]*bodyPlanNode, 0),
	}
}
func (s bodyPlanNodeStack) Empty() bool          { return len(s.v) == 0 }
func (s bodyPlanNodeStack) Peek() *bodyPlanNode  { return s.v[len(s.v)-1] }
func (s bodyPlanNodeStack) Len() int             { return len(s.v) }
func (s *bodyPlanNodeStack) Put(n *bodyPlanNode) { s.v = append(s.v, n) }
func (s *bodyPlanNodeStack) Pop() *bodyPlanNode {
	d := s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]
	return d
}
