package process

import (
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
)

// A stochastic process can generate a sequence of random variables representing
// the change of state of some system over time. The particular variables and
// their distributions depend on the particular process.
type StochasticProcess interface {

	// Generate the next random variable from the process
	Sample() variable.RandomVariable

	// Generate the next n random variables from the process
	SampleN(n int) []variable.RandomVariable

	// Return factors relating the process parameters to the given sequence
	Factors(sequence []variable.RandomVariable) []factor.Factor
}

// A default implementation of SampleN() for a StochasticProcess
type DefProcessDistSampleN struct{ process StochasticProcess }

func (def DefProcessDistSampleN) SampleN(n int) []variable.RandomVariable {
	var vars []variable.RandomVariable
	for i := 0; i < n; i++ {
		vars = append(vars, def.process.Sample())
	}
	return vars
}
