package astar

import (
	"math"
)

type Node interface {
	Neighbours() []Node
	HeuristicCostEstimate(Node) float64
}

type marker struct {
	gScore, fScore float64
	cameFrom       *marker
	node           Node
}

type markers map[*marker]struct{}

func (ns markers) Add(node *marker) {
	ns[node] = struct{}{}
}

func (ns markers) Drop(node *marker) {
	delete(ns, node)
}

func (ns markers) Has(node *marker) bool {
	_, has := ns[node]
	return has
}

func Search(src, dst Node) []Node {

	inf := math.Inf(+1)
	nodes := map[Node]*marker{}

	node := func(n Node) *marker {
		if _, exists := nodes[n]; !exists {
			nodes[n] = &marker{
				gScore:   inf,
				fScore:   inf,
				cameFrom: nil,
				node:     n,
			}
		}
		return nodes[n]
	}

	closedSet := markers{}
	openSet := markers{}

	start := node(src)
	start.gScore = 0
	start.fScore = start.node.HeuristicCostEstimate(dst)

	openSet.Add(start)

	for len(openSet) > 0 {

		var current *marker

		for candidate, _ := range openSet {
			if current == nil || candidate.fScore < current.fScore {
				current = candidate
			}
		}

		if current.node == dst {

			var nodes []Node

			for current.cameFrom != nil {
				nodes = append(nodes, current.node)
				current = current.cameFrom
			}

			nodes = append(nodes, src)

			for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}

			return nodes

		} else {

			openSet.Drop(current)
			closedSet.Add(current)

			for _, n := range current.node.Neighbours() {
				neighbour := node(n)
				if !closedSet.Has(neighbour) {
					gScoreTentative := current.gScore + current.node.HeuristicCostEstimate(neighbour.node)

					if !openSet.Has(neighbour) || gScoreTentative < neighbour.gScore {
						neighbour.cameFrom = current
						neighbour.gScore = gScoreTentative
						neighbour.fScore = gScoreTentative + neighbour.node.HeuristicCostEstimate(dst)
						openSet.Add(neighbour)
					}
				}
			}
		}
	}

	return nil
}
