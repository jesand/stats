package factor

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/dist"
)

// A connecting node in a factor graph. A factor is a node with edges to
// random variable nodes, and which has a corresponding function to score the
// values of those random variables.
type Factor interface {

	// The adjacent random variables
	Adjacent() []RandomVariable

	// The factor's current score, based on the values of adjacent variables
	Score() float64
}

// Create a new factor which scores based on a probability distribution
func NewDistFactor(vars []RandomVariable, distr dist.Dist) *DistFactor {
	return &DistFactor{
		Vars: vars,
		Dist: distr,
	}
}

// A factor which scores variables based on a probability distribution
type DistFactor struct {
	Vars []RandomVariable
	Dist dist.Dist
}

// The adjacent random variables
func (factor DistFactor) Adjacent() []RandomVariable {
	return factor.Vars
}

// The log probability of the variables given the parameters
func (factor DistFactor) Score() float64 {
	var (
		numVars   = factor.Dist.NumVars()
		numParams = factor.Dist.NumParams()
	)
	if len(factor.Vars) != numVars+numParams {
		panic(stats.ErrfFactorVarNum(numVars, numParams, len(factor.Vars)))
	}
	var (
		vars   = make([]float64, numVars)
		params = make([]float64, numParams)
	)
	for i, rv := range factor.Vars {
		if i < len(vars) {
			vars[i] = rv.Val()
		} else {
			params[i-len(vars)] = rv.Val()
		}
	}
	return factor.Dist.Score(vars, params)
}
