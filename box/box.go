package box

import (
	"image"
)

const (
	Left = iota
	Right
	Top
	Bottom
)

type Box struct {
	X int
	Y int
	W int
	H int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (self Box) Rect() image.Rectangle {
	return image.Rect(self.X, self.Y, self.X+self.W, self.Y+self.H)
}

// Extract X,Y,W,H
func (self Box) Unpack() (int, int, int, int) {
	return self.X, self.Y, self.W, self.H
}

// Alter X,Y
func (self Box) Translate(x, y int) Box {
	self.X += x
	self.Y += y
	return self
}

// Alter W,H
func (self Box) Extend(w, h int) Box {
	self.W += w
	self.H += h
	return self
}

// Grow in place (centroid unchanged)
func (self Box) Grow(x, y int) Box {
	self.X -= x
	self.Y -= y
	self.W += x * 2
	self.H += y * 2
	return self
}

// Vertical 50:50 split returning left and right boxes
func (self Box) Vsplit() (Box, Box) {
	x, y, w, h := self.Unpack()
	return Box{x, y, w / 2, h}, Box{x + w/2, y, w / 2, h}
}

// Horizontal 50:50 split returning above and bolow boxes
func (self Box) Hsplit() (Box, Box) {
	x, y, w, h := self.Unpack()
	return Box{x, y, w, h / 2}, Box{x, y + h/2, w, h / 2}
}

// Are boxes equal? Go can also just == the structs
func (self Box) Equal(b Box) bool {
	return self.Rect().Eq(b.Rect())
}

// Do boxes intersect in any direction?
func (self Box) Intersects(b Box) bool {
	return !self.Rect().Intersect(b.Rect()).Empty()
}

// Are boxes adjacent with one side touching?
func (self Box) Adjacent(b Box) bool {
	a := self
	if a.Intersects(b) {
		return false
	}
	ax, ay, aw, ah := a.Unpack()
	bx, by, bw, bh := b.Unpack()

	alignX := ax == bx+bw || bx == ax+aw
	alignY := ay == by+bh || by == ay+ah

	within := func(a, b, c int) bool {
		return c >= a && c < b
	}

	overlapX := within(ax, ax+aw, bx) || within(ax, ax+aw, bx+bw-1) || within(bx, bx+bw, ax) || within(bx, bx+bw, ax+aw-1)
	overlapY := within(ay, ay+ah, by) || within(ay, ay+ah, by+bh-1) || within(by, by+bh, ay) || within(by, by+bh, ay+ah-1)

	return (alignX && overlapY) || (alignY && overlapX)
}

// Is a box entirely contained within another?
func (self Box) Contains(b Box) bool {
	return b.Rect().In(self.Rect())
}

// Bounding box
func (self Box) Union(b Box) Box {
	r := self.Rect().Union(b.Rect())
	return Box{r.Min.X, r.Min.Y, r.Dx(), r.Dy()}
}

// X,Y of center point
func (self Box) Centroid() (int, int) {
	x, y, w, h := self.Unpack()
	return int(float64(x) + float64(w)/2), int(float64(y) + float64(h)/2)
}

// Split a box into rows and cols returning one cell
func (self Box) Cell(cols, rows, col, row int) Box {
	return Box{self.X + self.W*col, self.Y + self.H*row, self.W / cols, self.H / rows}
}

// Move a box side to side within another
func (self Box) Float(b Box, dir int) Box {
	switch dir {
	case Left:
		b.X = self.X
	case Right:
		b.X = self.X + self.W - b.W
	case Top:
		b.Y = self.Y
	case Bottom:
		b.Y = self.Y + self.H - b.H
	}
	return b
}
