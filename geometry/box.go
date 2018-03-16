package geometry

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

func (box Box) Rect() image.Rectangle {
	return image.Rect(box.X, box.Y, box.X+box.W, box.Y+box.H)
}

// Extract X,Y,W,H
func (box Box) Unpack() (int, int, int, int) {
	return box.X, box.Y, box.W, box.H
}

func (box Box) Point() Point {
	return Point{box.X, box.Y}
}

// Alter X,Y
func (box Box) Translate(x, y int) Box {
	box.X += x
	box.Y += y
	return box
}

// Alter W,H
func (box Box) Extend(w, h int) Box {
	box.W += w
	box.H += h
	return box
}

// Grow in place (centroid unchanged)
func (box Box) Grow(x, y int) Box {
	box.X -= x
	box.Y -= y
	box.W += x * 2
	box.H += y * 2
	return box
}

// Vertical 50:50 split returning left and right boxes
func (box Box) Vsplit() (Box, Box) {
	x, y, w, h := box.Unpack()
	return Box{x, y, w / 2, h}, Box{x + w/2, y, w / 2, h}
}

// Horizontal 50:50 split returning above and bolow boxes
func (box Box) Hsplit() (Box, Box) {
	x, y, w, h := box.Unpack()
	return Box{x, y, w, h / 2}, Box{x, y + h/2, w, h / 2}
}

// Are boxes equal? Go can also just == the structs
func (box Box) Equal(b Box) bool {
	return box.Rect().Eq(b.Rect())
}

// Do boxes intersect in any direction?
func (box Box) Intersects(b Box) bool {
	return !box.Rect().Intersect(b.Rect()).Empty()
}

// Are boxes adjacent with one side touching?
func (box Box) Adjacent(b Box) bool {
	a := box
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
func (box Box) Contains(b Box) bool {
	return b.Rect().In(box.Rect())
}

// Bounding box
func (box Box) Union(b Box) Box {
	r := box.Rect().Union(b.Rect())
	return Box{r.Min.X, r.Min.Y, r.Dx(), r.Dy()}
}

// X,Y of center point
func (box Box) Centroid() (int, int) {
	x, y, w, h := box.Unpack()
	return int(float64(x) + float64(w)/2), int(float64(y) + float64(h)/2)
}

// Split a box into rows and cols returning one cell
func (box Box) Cell(cols, rows, col, row int) Box {
	return Box{box.X + box.W*col, box.Y + box.H*row, box.W / cols, box.H / rows}
}

// Move a box side to side within another
func (box Box) Float(b Box, dir int) Box {
	switch dir {
	case Left:
		b.X = box.X
	case Right:
		b.X = box.X + box.W - b.W
	case Top:
		b.Y = box.Y
	case Bottom:
		b.Y = box.Y + box.H - b.H
	}
	return b
}
