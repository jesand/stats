package dist

import (
	"math"
	"math/rand"
)

// Represents a distribution over reals for a random variable
type RealDist interface {

	// The value of the CDF: Pr(X <= val) for random variable X over this space
	CDF(val float64) float64

	// The mean, or expected value, of the random variable
	Mean() float64

	// The mode of the random variable
	Mode() float64

	// The variance of the random variable
	Variance() float64
}

// Represents a continuous distribution over a subset of reals
type ContinuousDist interface {
	RealDist

	// Sample an outcome from the distribution
	Sample() Outcome

	// Sample a sequence of n outcomes from the distribution
	SampleN(n int) []Outcome

	// Return the corresponding sample space
	Space() DiscreteSpace

	// Return the probability of a given interval
	Prob(from, to float64) float64

	// Return the log probability (base 2) of a given interval
	LgProb(from, to float64) float64
}

// Represents a discrete distribution over a sample space
type DiscreteDist interface {

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
	panic(ErrNotNormalized)
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
