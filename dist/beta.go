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
		space: UnitIntervalSpace,
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

// Return a "score" (density or probability) for the given values
func (dist Beta) Score(vars, params []float64) float64 {
	return Beta{Alpha: params[0], Beta: params[1]}.PDF(vars[0])
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
	pdf := stat.Beta_PDF_At(dist.Alpha, dist.Beta, val)

	// Approximate difficult-to-calculate scores with a Normal distribution.
	// This approximation works well when the parameters are large.
	if dist.Alpha > 100 && dist.Beta > 100 &&
		(math.IsNaN(pdf) || math.IsInf(pdf, +1) || math.IsInf(pdf, -1)) {
		return Normal{
			Mu:    dist.Mean(),
			Sigma: math.Sqrt(dist.Variance()),
		}.PDF(val)
	}
	return pdf
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

// Return the Bayesian posterior using this Beta as a prior distribution, and
// having observed `pos` positive and `neg` negative outcomes.
func (dist Beta) Posterior(pos, neg float64) *Beta {
	return NewBetaDist(dist.Alpha+pos, dist.Beta+neg)
}

// Return the Beta which maximizes the probability of emitting the given
// sequence, based on a method of moments estimation.
// See: http://www.itl.nist.gov/div898/handbook/eda/section3/eda366h.htm
func (dist Beta) MaximizeByMoM(vals []float64) *Beta {
	var (
		mean     = Mean(vals)
		variance = Variance(vals)
		scale    = ((mean * (1 - mean)) / variance) - 1
		alpha    = mean * scale
		beta     = (1 - mean) * scale
	)
	return NewBetaDist(alpha, beta)
}
