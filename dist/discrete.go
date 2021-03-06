package dist

import (
	"github.com/jesand/stats"
)

// Make a new instance of DenseMutableDiscreteDist
func NewDenseMutableDiscreteDist(space DiscreteSpace) *DenseMutableDiscreteDist {
	dist := &DenseMutableDiscreteDist{
		space:   space,
		weights: make([]float64, space.Size()),
	}
	dist.DefDiscreteDistSample.dist = dist
	dist.DefDiscreteDistSampleN.dist = dist
	dist.DefDiscreteDistLgProb.dist = dist
	return dist
}

// A mutable discrete distribution which stores a dense probability vector
// for its outcomes.
type DenseMutableDiscreteDist struct {
	DefDiscreteDistSample
	DefDiscreteDistSampleN
	DefDiscreteDistLgProb

	// The sample space
	space DiscreteSpace

	// The probability mass for each outcome
	weights []float64

	// The total probability mass, used for normalization
	totalWeight float64
}

// Return a "score" (density or probability) for the given values
func (dist DenseMutableDiscreteDist) Score(vars, params []float64) float64 {
	outcome := dist.space.(DiscreteRealSpace).Outcome(vars[0])
	return params[int(outcome)]
}

// The number of random variables the distribution is over
func (dist DenseMutableDiscreteDist) NumVars() int {
	return 1
}

// The number of parameters in the distribution: the weights
func (dist DenseMutableDiscreteDist) NumParams() int {
	return len(dist.weights)
}

// Update the distribution parameters
func (dist *DenseMutableDiscreteDist) SetParams(vals []float64) {
	copy(dist.weights[:], vals[:])
}

// Return the sample space
func (dist DenseMutableDiscreteDist) Space() DiscreteSpace {
	return dist.space
}

// Return the probability of a given outcome
func (dist DenseMutableDiscreteDist) Prob(outcome Outcome) float64 {
	if dist.totalWeight != 1 {
		panic(stats.ErrNotNormalized)
	} else if int(outcome) < 0 || int(outcome) >= len(dist.weights) {
		panic(stats.ErrfNotInDomain(int(outcome)))
	}
	return dist.weights[int(outcome)]
}

// Set the probability of a particular outcome
func (dist *DenseMutableDiscreteDist) SetProb(outcome Outcome, prob float64) {
	if int(outcome) < 0 || int(outcome) >= len(dist.weights) {
		panic(stats.ErrfNotInDomain(int(outcome)))
	} else if prob < 0 || prob > 1 {
		panic(stats.ErrfInvalidProb(prob))
	}
	dist.totalWeight -= dist.weights[outcome]
	dist.weights[outcome] = prob
	dist.totalWeight += dist.weights[outcome]
}

// Set the unnormalized measure for a particular outcome. It is
// up to the particular distribution to normalize these weights.
func (dist *DenseMutableDiscreteDist) SetWeight(outcome Outcome, weight float64) {
	dist.SetProb(outcome, weight)
}

// Set all probabilities to zero
func (dist *DenseMutableDiscreteDist) Reset() {
	for i := range dist.weights {
		dist.weights[i] = 0
	}
	dist.totalWeight = 0
}

// Normalize all weights, assuming 0 weight for outcomes not assigned with
// SetWeight() since the last call to Normalize().
func (dist *DenseMutableDiscreteDist) Normalize() {
	if dist.totalWeight == 0 {
		panic(stats.ErrZeroProb)
	} else if dist.totalWeight != 1 {
		for i, w := range dist.weights {
			dist.weights[i] = w / dist.totalWeight
		}
		dist.totalWeight = 1
	}
}

// Normalize all weights, assigning `rest` weight uniformly to all outcomes
// currently assigned zero weight.
func (dist *DenseMutableDiscreteDist) NormalizeWithExtra(rest float64) {
	if rest != 0 {
		var numZeros float64
		for _, w := range dist.weights {
			if w == 0 {
				numZeros++
			}
		}
		if numZeros > 0 {
			var weight = rest / numZeros
			for i, w := range dist.weights {
				if w == 0 {
					dist.weights[i] = weight
				}
			}
		}
	}
	dist.Normalize()
}
