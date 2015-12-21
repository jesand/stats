package model

import (
	"fmt"
	"math/rand"
)

// Generates a random bipartite graph with the specified node degrees. The graph
// is selected approximately uniformly at random from the space of all graphs
// with the specified node degrees. Implements the algorithm in:
//
// M. Bayati, J. H. Kim, and A. Saberi, "A Sequential Algorithm for Generating
// Random Graphs," Algorithmica, vol. 58, no. 4, pp. 860–910, 2010.
//
// If successful, returns a list of edges. Each edge contains an index from the
// left nodes and an index from the right nodes.
func RandomBipartiteGraph(leftDegrees, rightDegrees []int) (edges [][2]int, err error) {
	var (
		leftNeeded            = make([]int, len(leftDegrees))
		rightNeeded           = make([]int, len(rightDegrees))
		totalLeft, totalRight int
	)
	copy(leftNeeded[:], leftDegrees[:])
	copy(rightNeeded[:], rightDegrees[:])
	for _, d := range leftDegrees {
		totalLeft += d
	}
	for _, d := range rightDegrees {
		totalRight += d
	}
	if totalLeft == 0 && totalRight == 0 {
		return nil, fmt.Errorf("Total node degree is zero")
	} else if totalLeft != totalRight {
		return nil, fmt.Errorf("Total left degree %d != total right degree %d",
			totalLeft, totalRight)
	}
	var hasEdge = make([]bool, len(leftNeeded)*len(rightNeeded))
	edges = make([][2]int, totalLeft)
	for edgeNum := 0; edgeNum < totalLeft; edgeNum++ {

		// Find the total probability of unselected edges
		var totalWeight float64
		for l, ld := range leftDegrees {
			for r, rd := range rightDegrees {
				idx := l*len(rightDegrees) + r
				if l != r && !hasEdge[idx] {
					w := float64(leftNeeded[l]*rightNeeded[r]) *
						(1 - float64(ld*rd)/float64(4*totalLeft))
					totalWeight += w
				}
			}
		}
		if totalWeight == 0 {
			return nil, fmt.Errorf("Could not find a random graph of the given degree")
		}

		// Select an edge
		var (
			remaining = rand.Float64() * totalWeight
			found     = false
		)
		for l, ld := range leftDegrees {
			for r, rd := range rightDegrees {
				idx := l*len(rightDegrees) + r
				if l != r && !hasEdge[idx] {
					w := float64(leftNeeded[l]*rightNeeded[r]) *
						(1 - float64(ld*rd)/float64(4*totalLeft))
					remaining -= w
					if remaining <= 0 {
						edges[edgeNum] = [2]int{l, r}
						leftNeeded[l]--
						rightNeeded[r]--
						hasEdge[idx] = true
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}
	}
	return
}
