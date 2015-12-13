package process

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
)

// Create a new IIDProcess based on the given distribution
func NewIIDProcess(params []variable.RandomVariable, dist dist.Dist) *IIDProcess {
	return &IIDProcess{
		Params: params,
		Dist:   dist,
	}
}

// A StochasticProcess which draws iid variables from an underlying distribution
type IIDProcess struct {

	// Distribution parameters
	Params []variable.RandomVariable

	// The distribution
	Dist dist.Dist
}

// Generate the next random variable from the process
func (process IIDProcess) Sample() variable.RandomVariable {
	if source, ok := process.Dist.(dist.DiscreteDist); ok {
		space := source.Space().(dist.DiscreteRealSpace)
		return variable.NewDiscreteRV(source.Sample(), space)
	} else if source, ok := process.Dist.(dist.ContinuousDist); ok {
		return variable.NewContinuousRV(source.Sample(), source.Space())
	} else {
		panic(stats.ErrfUnsupportedDist(process.Dist))
	}
}

// Generate the next n random variables from the process
func (process IIDProcess) SampleN(n int) (rvs []variable.RandomVariable) {
	if source, ok := process.Dist.(dist.DiscreteRealDist); ok {
		space := source.Space().(dist.DiscreteRealSpace)
		for _, v := range source.SampleN(n) {
			rvs = append(rvs, variable.NewDiscreteRV(v, space))
		}
		return rvs
	} else if source, ok := process.Dist.(dist.ContinuousDist); ok {
		space := source.Space()
		for _, v := range source.SampleN(n) {
			rvs = append(rvs, variable.NewContinuousRV(v, space))
		}
		return rvs
	} else {
		panic(stats.ErrfUnsupportedDist(process.Dist))
	}
}

// Return factors relating the process parameters to the given sequence
func (process IIDProcess) Factors(sequence []variable.RandomVariable) (
	factors []factor.Factor) {

	for i := range sequence {
		factors = append(factors, factor.NewDistFactor(
			append(sequence[i:i+1], process.Params...),
			process.Dist,
		))
	}
	return
}
