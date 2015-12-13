package dist

import (
	"github.com/jesand/stats"
	"math"
	"math/rand"
)

// Represents a probability distribution
type Dist interface {

	// Return a "score" (log density or log mass) for the given values
	Score(vars, params []float64) float64

	// The number of random variables the distribution is over
	NumVars() int

	// The number of parameters in the distribution
	NumParams() int

	// Update the distribution parameters
	SetParams(vals []float64)
}

// Represents a distribution over reals for a random variable
type RealDist interface {

	// The mean, or expected value, of the random variable
	Mean() float64

	// The mode of the random variable
	Mode() float64

	// The variance of the random variable
	Variance() float64
}

// Represents a continuous distribution over a subset of reals
type ContinuousDist interface {
	Dist
	RealDist

	// Sample an outcome from the distribution
	Sample() float64

	// Sample a sequence of n outcomes from the distribution
	SampleN(n int) []float64

	// Return the corresponding sample space
	Space() RealSpace

	// The value of the CDF: Pr(X <= val) for random variable X over this space
	CDF(val float64) float64

	// Return the density at a given value
	PDF(val float64) float64

	// Return the probability of a given interval
	Prob(from, to float64) float64

	// Return the log probability (base 2) of a given interval
	LgProb(from, to float64) float64
}

// Represents a discrete distribution over a sample space
type DiscreteDist interface {
	Dist

	// Sample an outcome from the distribution
	Sample() Outcome

	// Sample a sequence of n outcomes from the distribution
	SampleN(n int) []Outcome

	// Return the corresponding sample space
	Space() DiscreteSpace

	// Return the probability of a given outcome
	Prob(outcome Outcome) float64

	// Return the log probability (base 2) of a given outcome
	LgProb(outcome Outcome) float64
}

// A discrete distribution whose underlying probability measure can change
type MutableDiscreteDist interface {

	// A mutable dist is a dist
	DiscreteDist

	// Set all probabilities to zero
	Reset()

	// Set the probability of a particular outcome
	SetProb(outcome Outcome, prob float64)

	// Set the unnormalized measure for a particular outcome. It is
	// up to the particular distribution to normalize these weights.
	SetWeight(outcome Outcome, weight float64)

	// Normalize all weights, assuming 0 weight for outcomes not assigned with
	// SetWeight() since the last call to Normalize().
	Normalize()

	// Normalize all weights, assigning `rest` weight uniformly to all outcomes
	// currently assigned zero weight.
	NormalizeWithExtra(rest float64)
}

// A default implementation of Sample() for a DiscreteDist
type DefDiscreteDistSample struct{ dist DiscreteDist }

func (dist DefDiscreteDistSample) Sample() Outcome {
	var remaining = rand.Float64()
	for i := Outcome(0); int(i) < dist.dist.Space().Size(); i++ {
		remaining -= dist.dist.Prob(i)
		if remaining <= 0 {
			return i
		}
	}
	panic(stats.ErrNotNormalized)
}

// A default implementation of SampleN() for a DiscreteDist
type DefDiscreteDistSampleN struct{ dist DiscreteDist }

func (dist DefDiscreteDistSampleN) SampleN(n int) []Outcome {
	var outcomes []Outcome
	for i := 0; i < n; i++ {
		outcomes = append(outcomes, dist.dist.Sample())
	}
	return outcomes
}

// A default implementation of LgProb() for a DiscreteDist
type DefDiscreteDistLgProb struct{ dist DiscreteDist }

// Return the log probability (base 2) of a given outcome
func (dist DefDiscreteDistLgProb) LgProb(outcome Outcome) float64 {
	return math.Log2(dist.dist.Prob(outcome))
}

// A default implementation of SampleN() for a ContinuousDist
type DefContinuousDistSampleN struct{ dist ContinuousDist }

func (dist DefContinuousDistSampleN) SampleN(n int) []float64 {
	var outcomes []float64
	for i := 0; i < n; i++ {
		outcomes = append(outcomes, dist.dist.Sample())
	}
	return outcomes
}

// A default implementation of Prob() for a ContinuousDist
type DefContinuousDistProb struct{ dist ContinuousDist }

// Return the log probability (base 2) of a given outcome
func (dist DefContinuousDistProb) Prob(from, to float64) float64 {
	return dist.dist.CDF(to) - dist.dist.CDF(from)
}

// A default implementation of LgProb() for a ContinuousDist
type DefContinuousDistLgProb struct{ dist ContinuousDist }

// Return the log probability (base 2) of a given outcome
func (dist DefContinuousDistLgProb) LgProb(from, to float64) float64 {
	return math.Log2(dist.dist.Prob(from, to))
}
