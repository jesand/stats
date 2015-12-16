package factor

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/variable"
)

// A connecting node in a factor graph. A factor is a node with edges to
// random variable nodes, and which has a corresponding function to score the
// values of those random variables.
type Factor interface {

	// The adjacent random variables
	Adjacent() []variable.RandomVariable

	// The factor's current score, based on the values of adjacent variables
	Score() float64
}

// Create a new factor which scores based on a probability distribution.
// The variables are split into "variables" and "parameters" using the
// distribution's NumVars() and NumParams() values.
func NewDistFactor(vars []variable.RandomVariable, distr dist.Dist) *DistFactor {
	return &DistFactor{
		Vars: vars,
		Dist: distr,
	}
}

// A factor which scores variables based on a probability distribution
type DistFactor struct {
	Vars []variable.RandomVariable
	Dist dist.Dist
}

// The adjacent random variables
func (factor DistFactor) Adjacent() []variable.RandomVariable {
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

// Create a new factor which always returns the same score
func NewConstFactor(vars []variable.RandomVariable, value float64) *ConstFactor {
	return &ConstFactor{
		Vars:  vars,
		Value: value,
	}
}

// A Factor which always returns the same value
type ConstFactor struct {
	Vars  []variable.RandomVariable
	Value float64
}

// The adjacent random variables
func (factor ConstFactor) Adjacent() []variable.RandomVariable {
	return factor.Vars
}

// The log probability of the variables given the parameters
func (factor ConstFactor) Score() float64 {
	return factor.Value
}
