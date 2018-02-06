package box

import (
	"testing"
)

func TestUnpack(test *testing.T) {
	a := Box{1, 2, 3, 4}
	x, y, w, h := a.Unpack()
	if x != 1 || y != 2 || w != 3 || h != 4 {
		test.Errorf("Unpack()")
	}
}

func TestTranslate(test *testing.T) {
	a := Box{1, 2, 3, 4}
	a = a.Translate(10, 10)
	if a.X != 11 || a.Y != 12 {
		test.Errorf("Translate()")
	}
}

func TestExtend(test *testing.T) {
	a := Box{1, 2, 3, 4}
	a = a.Extend(10, 10)
	if a.W != 13 || a.H != 14 {
		test.Errorf("Extend()")
	}
}

func TestGrow(test *testing.T) {
	a := Box{1, 2, 3, 4}
	a = a.Grow(10, 10)
	a1 := Box{-9, -8, 23, 24}
	if a != a1 {
		test.Errorf("Grow()")
	}
}

func TestVsplit(test *testing.T) {
	a := Box{10, 10, 10, 10}
	b, c := a.Vsplit()
	b1 := Box{10, 10, 5, 10}
	c1 := Box{15, 10, 5, 10}
	if b != b1 || c != c1 {
		test.Errorf("Vsplit()")
	}
}

func TestHsplit(test *testing.T) {
	a := Box{10, 10, 10, 10}
	b, c := a.Hsplit()
	b1 := Box{10, 10, 10, 5}
	c1 := Box{10, 15, 10, 5}
	if b != b1 || c != c1 {
		test.Errorf("Hsplit()")
	}
}

func TestEqual(test *testing.T) {
	a := Box{10, 10, 10, 10}
	b := Box{10, 10, 10, 10}
	c := Box{1, 2, 3, 4}
	if !a.Equal(b) || b.Equal(c) != false {
		test.Errorf("Hsplit()")
	}
}

func TestIntersects(test *testing.T) {
	a := Box{10, 10, 10, 10}
	if !a.Intersects(Box{15, 15, 5, 5}) {
		test.Errorf("Intersects() -- internal")
	}
	if !a.Intersects(Box{15, 5, 5, 10}) {
		test.Errorf("Intersects() -- top")
	}
	if !a.Intersects(Box{15, 15, 5, 10}) {
		test.Errorf("Intersects() -- bottom")
	}
	if !a.Intersects(Box{5, 15, 10, 5}) {
		test.Errorf("Intersects() -- left")
	}
	if !a.Intersects(Box{15, 15, 10, 5}) {
		test.Errorf("Intersects() -- right")
	}
	if a.Intersects(Box{20, 20, 10, 10}) {
		test.Errorf("Intersects() -- external")
	}
}

func TestAdjacent(test *testing.T) {
	a := Box{10, 10, 10, 10}
	if !a.Adjacent(Box{10, 5, 5, 5}) {
		test.Errorf("Adjacent() -- top")
	}
	if !a.Adjacent(Box{10, 20, 5, 5}) {
		test.Errorf("Adjacent() -- bottom")
	}
	if !a.Adjacent(Box{0, 10, 10, 5}) {
		test.Errorf("Adjacent() -- left")
	}
	if !a.Adjacent(Box{20, 10, 10, 10}) {
		test.Errorf("Adjacent() -- right")
	}
	if a.Adjacent(Box{21, 10, 10, 10}) {
		test.Errorf("Adjacent() -- right offset")
	}
}

func TestContains(test *testing.T) {
	a := Box{10, 10, 10, 10}
	if !a.Contains(Box{10, 10, 5, 5}) {
		test.Errorf("Contains() -- internal")
	}
	if a.Contains(Box{10, 10, 50, 5}) {
		test.Errorf("Contains() -- external")
	}
}

func TestUnion(test *testing.T) {
	a := Box{10, 10, 10, 10}
	b := Box{15, 10, 10, 5}
	if a.Union(b) != (Box{10, 10, 15, 10}) {
		test.Errorf("Union()")
	}
}

func TestCentroid(test *testing.T) {
	a := Box{10, 10, 10, 10}
	x, y := a.Centroid()
	if x != 15 || y != 15 {
		test.Errorf("Centroid()")
	}
}

func TestCell(test *testing.T) {
	a := Box{10, 10, 10, 10}
	if a.Cell(2, 2, 0, 0) != (Box{10, 10, 5, 5}) {
		test.Errorf("Cell()")
	}
}

func TestFloat(test *testing.T) {
	a := Box{10, 10, 10, 10}
	if a.Float(Box{15, 15, 1, 1}, Left) != (Box{10, 15, 1, 1}) {
		test.Errorf("Float()")
	}
}
