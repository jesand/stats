package factor

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/variable"
)

// Create a new factor graph
func NewFactorGraph() *FactorGraph {
	return &FactorGraph{
		varIds: make(map[variable.RandomVariable]int),
	}
}

// A FactorGraph is a bipartite graph between random variables and factors.
// The joint probability distribution over all random variables is the product
// of all factors, calculated over their adjacent random variables.
type FactorGraph struct {
	Factors   []Factor
	Variables []factorGraphVar
	varIds    map[variable.RandomVariable]int
}

type factorGraphVar struct {
	Variable variable.RandomVariable
	Factors  []Factor
}

// Add a factor and its adjacent random variables to the graph
func (graph *FactorGraph) AddFactor(factor Factor) {
	graph.Factors = append(graph.Factors, factor)
	for _, v := range factor.Adjacent() {
		idx, ok := graph.varIds[v]
		if ok {
			graph.Variables[idx].Factors = append(graph.Variables[idx].Factors,
				factor)
		} else {
			idx = len(graph.Variables)
			graph.Variables = append(graph.Variables,
				factorGraphVar{v, []Factor{factor}})
			graph.varIds[v] = idx
		}
	}
}

// Get variables adjacent to a factor
func (graph FactorGraph) AdjToFactor(factor Factor) []variable.RandomVariable {
	return factor.Adjacent()
}

// Get factors adjacent to a variable
func (graph FactorGraph) AdjToVariable(v variable.RandomVariable) []Factor {
	if idx, ok := graph.varIds[v]; !ok {
		panic(stats.Errorf("Random variable %#v not in factor graph", v))
	} else {
		return graph.Variables[idx].Factors
	}
}
