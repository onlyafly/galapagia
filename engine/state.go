package engine

type State struct {
}

func NewState() *State {
	return &State{}
}

func CurrentGrid() [][]int {
	g := make([][]int, 100)
	for i, _ := range g {
		g[i] = make([]int, 100)
	}
	g[1][20] = 1
	g[1][21] = 2
	return g
}
