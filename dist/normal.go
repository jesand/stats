package dist

import (
	"math"
	"math/rand"
)

// Produce a new Normal distribution
func NewNormalDist(mean, stdDev float64) *Normal {
	dist := &Normal{
		Mu:    mean,
		Sigma: stdDev,
		space: AllRealSpace,
	}
	dist.DefContinuousDistSampleN.dist = dist
	dist.DefContinuousDistProb.dist = dist
	dist.DefContinuousDistLgProb.dist = dist
	return dist
}

// Produce a new standard Normal distribution
func NewStandardNormalDist() *Normal {
	return NewNormalDist(0, 1)
}

// A Normal distribution. See: https://en.wikipedia.org/wiki/Normal_distribution
type Normal struct {

	// The distribution parameters
	Mu, Sigma float64

	// The space
	space RealSpace

	DefContinuousDistSampleN
	DefContinuousDistProb
	DefContinuousDistLgProb
}

// Return the corresponding sample space
func (dist Normal) Space() RealSpace {
	return dist.space
}

// Return a "score" (density or probability) for the given values
func (dist Normal) Score(vars, params []float64) float64 {
	return Normal{Mu: params[0], Sigma: params[1]}.PDF(vars[0])
}

// The number of random variables the distribution is over
func (dist Normal) NumVars() int {
	return 1
}

// The number of parameters in the distribution
func (dist Normal) NumParams() int {
	return 2
}

// Update the distribution parameters
func (dist *Normal) SetParams(vals []float64) {
	dist.Mu, dist.Sigma = vals[0], vals[1]
}

// Return the density at a given value
func (dist Normal) PDF(val float64) float64 {
	return math.Exp(-(math.Pow(val-dist.Mu, 2))/
		(2*dist.Sigma*dist.Sigma)) /
		(dist.Sigma * math.Sqrt2 * math.SqrtPi)
}

// The value of the CDF: Pr(X <= val) for random variable X over this space
func (dist Normal) CDF(val float64) float64 {
	return (1 + math.Erf((val-dist.Mu)/(dist.Sigma*math.Sqrt2))) / 2
}

// The mean, or expected value, of the random variable
func (dist Normal) Mean() float64 {
	return dist.Mu
}

// The mode of the random variable
func (dist Normal) Mode() float64 {
	return dist.Mu
}

// The variance of the random variable
func (dist Normal) Variance() float64 {
	return dist.Sigma * dist.Sigma
}

// Sample an outcome from the distribution
func (dist Normal) Sample() float64 {
	return dist.Mu + rand.NormFloat64()/dist.Sigma
}
