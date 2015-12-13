package dist

// Create a new Bernoulli distribution
func NewBernoulliDist(bias float64) *BernoulliDist {
	dist := &BernoulliDist{
		DenseMutableDiscreteDist: NewDenseMutableDiscreteDist(BooleanSpace{}),
	}
	dist.SetBias(bias)
	return dist
}

// A Bernoulli distribution.
// See: https://en.wikipedia.org/wiki/Bernoulli_distribution
type BernoulliDist struct {
	*DenseMutableDiscreteDist
}

// Return a "score" (log density or log mass) for the given values
func (dist BernoulliDist) Score(vars, params []float64) float64 {
	return dist.DenseMutableDiscreteDist.Score(vars, []float64{0, params[0]})
}

// The number of parameters in the distribution: the weights
func (dist BernoulliDist) NumParams() int {
	return 1
}

// Update the distribution parameters
func (dist *BernoulliDist) SetParams(vals []float64) {
	dist.SetBias(vals[0])
}

// Return the space as a BooleanSpace
func (dist *BernoulliDist) BSpace() BooleanSpace {
	return dist.space.(BooleanSpace)
}

// Set the bias of the distribution
func (dist *BernoulliDist) SetBias(bias float64) {
	dist.SetProb(0, 1-bias)
	dist.SetProb(1, bias)
	dist.Normalize() // just in case, to avoid rounding errors
}

// The value of the CDF: Pr(X <= val) for random variable X over this space
func (dist BernoulliDist) CDF(val float64) float64 {
	if val < 0 {
		return 0
	} else if val < 1 {
		return dist.Prob(0)
	} else {
		return 1
	}
}

// The mean, or expected value, of the random variable
func (dist BernoulliDist) Mean() float64 {
	return dist.Prob(1)
}

// The mode of the random variable
func (dist BernoulliDist) Mode() float64 {
	if dist.Prob(0) > dist.Prob(1) {
		return dist.BSpace().F64Value(0)
	} else {
		return dist.BSpace().F64Value(1)
	}
}

// The variance of the random variable
func (dist BernoulliDist) Variance() float64 {
	return dist.Prob(0) * dist.Prob(1)
}
