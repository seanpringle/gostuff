package astar

import (
	"math"
	"testing"
)

type Maze [10][10]int

type Cell struct {
	X, Y int
	maze Maze
}

func (c Cell) Neighbours() []Node {
	var cells []Node
	tryXY := func(x, y int) {
		cn := Cell{X: x, Y: y, maze: c.maze}
		if c.maze.Traversable(cn) {
			cells = append(cells, cn)
		}
	}
	tryXY(c.X-1, c.Y-1)
	tryXY(c.X-1, c.Y+0)
	tryXY(c.X-1, c.Y+1)

	tryXY(c.X+0, c.Y-1)
	tryXY(c.X+0, c.Y+1)

	tryXY(c.X+1, c.Y-1)
	tryXY(c.X+1, c.Y+0)
	tryXY(c.X+1, c.Y+1)
	return cells
}

func (c Cell) HeuristicCostEstimate(t Node) float64 {
	dx := math.Abs(float64(t.(Cell).X - c.X))
	dy := math.Abs(float64(t.(Cell).Y - c.Y))
	return math.Sqrt(dx*dx + dy*dy)
}

func (m Maze) Traversable(c Node) (ok bool) {
	defer func() { recover() }()
	ok = m[c.(Cell).X][c.(Cell).Y] == 0
	return
}

func (m Maze) Validate(path []Node) bool {
	for _, p := range path {
		if !m.Traversable(p) {
			return false
		}
	}
	return true
}

func (m Maze) Draw(test *testing.T, path []Node) {
	dump := ""
	for x := 0; x < 10; x++ {
		dump += "\n"
		for y := 0; y < 10; y++ {
			wall := m[x][y] == 1
			used := false
			for _, p := range path {
				if p == (Cell{x, y, m}) {
					used = true
				}
			}
			if used && wall {
				dump += "! "
				continue
			}
			if used {
				dump += "* "
				continue
			}
			if wall {
				dump += "x "
				continue
			}
			dump += ". "
		}
	}
	test.Logf("%s", dump)
}

func TestSearch(test *testing.T) {
	maze := Maze{}
	maze[5][2] = 1
	maze[5][3] = 1
	maze[5][4] = 1
	maze[5][5] = 1
	maze[5][6] = 1
	maze[5][7] = 1
	path := Search(Cell{2, 2, maze}, Cell{8, 8, maze})
	if !maze.Validate(path) {
		test.Errorf("path crosses walls")
	}
	maze.Draw(test, path)
	path = Search(Cell{2, 2, maze}, Cell{5, 5, maze})
	if path != nil {
		test.Errorf("impossible path")
	}
	maze.Draw(test, path)
}
