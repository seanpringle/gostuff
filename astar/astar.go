package astar

import (
	"math"
)

type Point struct {
	X, Y int
}

type Pathable interface {
	Diagonal() bool
	Traversable(Point) bool
	HeuristicCostEstimate(Point, Point) float64
}

type navNode struct {
	pt             Point
	gScore, fScore float64
	cameFrom       *navNode
}

type navNodeSet map[*navNode]struct{}

func (ns navNodeSet) Add(node *navNode) {
	ns[node] = struct{}{}
}

func (ns navNodeSet) Drop(node *navNode) {
	delete(ns, node)
}

func (ns navNodeSet) Has(node *navNode) bool {
	_, has := ns[node]
	return has
}

func Search(cb Pathable, src, dst Point) []Point {

	diagonal := cb.Diagonal()

	inf := math.Inf(+1)
	nodes := make(map[Point]*navNode)

	node := func(pt Point) *navNode {
		if _, exists := nodes[pt]; !exists {
			nodes[pt] = &navNode{
				pt:       pt,
				gScore:   inf,
				fScore:   inf,
				cameFrom: nil,
			}
		}
		return nodes[pt]
	}

	closedSet := navNodeSet{}
	openSet := navNodeSet{}

	start := node(src)
	start.gScore = 0
	start.fScore = cb.HeuristicCostEstimate(src, dst)

	openSet.Add(start)

	var path []*navNode

	for path == nil && len(openSet) > 0 {

		var current *navNode

		for candidate, _ := range openSet {
			if current == nil || candidate.fScore < current.fScore {
				current = candidate
			}
		}

		if current.pt == dst {

			path = []*navNode{current}
			for current.cameFrom != nil {
				current = current.cameFrom
				path = append(path, current)
			}
			for i := 0; i < len(path)/2; i++ {
				tmp := path[i]
				path[i] = path[len(path)-i-1]
				path[len(path)-i-1] = tmp
			}

		} else {

			openSet.Drop(current)
			closedSet.Add(current)

			neighborCheck := func(pt Point) {
				neighbor := node(pt)
				// neighbor is obstacle?
				if !closedSet.Has(neighbor) && neighbor.pt != dst && !cb.Traversable(neighbor.pt) {
					closedSet.Add(neighbor)
				}
				if !closedSet.Has(neighbor) {
					gScoreTentative := current.gScore + cb.HeuristicCostEstimate(current.pt, neighbor.pt)

					if !openSet.Has(neighbor) || gScoreTentative < neighbor.gScore {
						neighbor.cameFrom = current
						neighbor.gScore = gScoreTentative
						neighbor.fScore = gScoreTentative + cb.HeuristicCostEstimate(neighbor.pt, dst)
						openSet.Add(neighbor)
					}
				}
			}

			if diagonal {
				neighborCheck(Point{current.pt.X - 1, current.pt.Y - 1})
				neighborCheck(Point{current.pt.X - 1, current.pt.Y + 1})
				neighborCheck(Point{current.pt.X + 1, current.pt.Y - 1})
				neighborCheck(Point{current.pt.X + 1, current.pt.Y + 1})
			}

			neighborCheck(Point{current.pt.X - 1, current.pt.Y + 0})
			neighborCheck(Point{current.pt.X + 0, current.pt.Y - 1})
			neighborCheck(Point{current.pt.X + 0, current.pt.Y + 1})
			neighborCheck(Point{current.pt.X + 1, current.pt.Y + 0})
		}
	}

	var points []Point

	for _, node := range path {
		points = append(points, node.pt)
	}

	return points
}
