package engine

import (
	"testing"
)

func Test_restrictToGrid(t *testing.T) {
	x, y := restrictToGrid(Size{10, 10}, 9, 9, 1, 1)
	if x != 9 || y != 9 {
		t.Error("Expected 9,9", x, y)
	}

	x, y = restrictToGrid(Size{10, 10}, 9, 9, 2, 2)
	if x != 8 || y != 8 {
		t.Error("Expected 8,8", x, y)
	}

	x, y = restrictToGrid(Size{100, 100}, 99, 99, 3, 3)
	if x != 97 || y != 96 {
		t.Error("Expected 97,97", x, y)
	}
}
