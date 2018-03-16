package geometry

type Circle struct {
	X int
	Y int
	R int
}

func (c Circle) Point() Point {
	return Point{c.X, c.Y}
}

func (c Circle) Box() Box {
	return Box{X: c.X - c.R, Y: c.Y - c.R, W: c.R * 2, H: c.R * 2}
}
