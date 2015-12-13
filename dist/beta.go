package dist

import (
	"github.com/ematvey/gostat"
	"github.com/jesand/stats"
	"math"
)

// Produce a new Beta distribution
func NewBetaDist(alpha, beta float64) *Beta {
	dist := &Beta{
		Alpha: alpha,
		Beta:  beta,
		space: NewUnitIntervalSpace(),
	}
	dist.DefContinuousDistSampleN.dist = dist
	dist.DefContinuousDistProb.dist = dist
	dist.DefContinuousDistLgProb.dist = dist
	return dist
}

// A Beta distribution. See: https://en.wikipedia.org/wiki/Beta_distribution
type Beta struct {

	// The distribution parameters
	Alpha, Beta float64

	// The space
	space RealSpace

	DefContinuousDistSampleN
	DefContinuousDistProb
	DefContinuousDistLgProb
}

// Return the corresponding sample space
func (dist Beta) Space() RealSpace {
	return dist.space
}

// Return a "score" (log density or log mass) for the given values
func (dist Beta) Score(vars, params []float64) float64 {
	alpha, beta := dist.Alpha, dist.Beta
	dist.Alpha, dist.Beta = params[0], params[1]
	score := math.Log2(dist.PDF(vars[0]))
	dist.Alpha, dist.Beta = alpha, beta
	return score
}

// The number of random variables the distribution is over
func (dist Beta) NumVars() int {
	return 1
}

// The number of parameters in the distribution
func (dist Beta) NumParams() int {
	return 2
}

// Update the distribution parameters
func (dist *Beta) SetParams(vals []float64) {
	dist.Alpha, dist.Beta = vals[0], vals[1]
}

// Return the density at a given value
func (dist Beta) PDF(val float64) float64 {
	return stat.Beta_PDF_At(dist.Alpha, dist.Beta, val)
}

// The value of the CDF: Pr(X <= val) for random variable X over this space
func (dist Beta) CDF(val float64) float64 {
	return stat.Beta_CDF_At(dist.Alpha, dist.Beta, val)
}

// The mean, or expected value, of the random variable
func (dist Beta) Mean() float64 {
	return dist.Alpha / (dist.Alpha + dist.Beta)
}

// The mode of the random variable
func (dist Beta) Mode() float64 {
	if dist.Alpha > 1 && dist.Beta > 1 {
		a, b := dist.Alpha, dist.Beta
		return (a - 1) / (a + b - 2)
	}
	panic(stats.Errorf("Beta(%f, %f) has no mode", dist.Alpha, dist.Beta))
}

// The variance of the random variable
func (dist Beta) Variance() float64 {
	a, b := dist.Alpha, dist.Beta
	return (a * b) / ((a + b) * (a + b) * (a + b + 1))
}

// Sample an outcome from the distribution
func (dist Beta) Sample() float64 {
	var (
		x = randGamma(dist.Alpha, 1, 0)
		y = randGamma(dist.Beta, 1, 0)
	)
	return x / (x + y)
}
