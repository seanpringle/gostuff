package astar

import (
	"math"
	"testing"
)

type Maze [10][10]int

func (m Maze) Diagonal() bool {
	return true
}

func (m Maze) Traversable(p Point) (ok bool) {
	defer func() { recover() }()
	ok = m[p.X][p.Y] == 0
	return
}

func (m Maze) HeuristicCostEstimate(a, b Point) float64 {
	dx := math.Abs(float64(b.X - a.X))
	dy := math.Abs(float64(b.Y - a.Y))
	return math.Sqrt(dx*dx + dy*dy)
}

func (m Maze) Validate(path []Point) bool {
	for _, p := range path {
		if !m.Traversable(p) {
			return false
		}
	}
	return true
}

func (m Maze) Draw(test *testing.T, path []Point) {
	dump := ""
	for x := 0; x < 10; x++ {
		dump += "\n"
		for y := 0; y < 10; y++ {
			wall := m[x][y] == 1
			used := false
			for _, p := range path {
				if p == (Point{x, y}) {
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
	maze[5][8] = 1
	path := Search(maze, Point{1, 1}, Point{8, 8})
	maze.Draw(test, path)
	if !maze.Validate(path) {
		test.Errorf("path crosses walls")
	}
}
