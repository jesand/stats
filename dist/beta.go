package dist

import (
	"github.com/ematvey/gostat"
	"github.com/jesand/stats"
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
	alpha, beta := dist.Alpha, dist.Beta
	dist.Alpha, dist.Beta = params[0], params[1]
	score := dist.PDF(vars[0])
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

// Return the Bayesian posterior using this Beta as a prior distribution, and
// having observed `pos` positive and `neg` negative outcomes.
func (dist Beta) Posterior(pos, neg float64) *Beta {
	return NewBetaDist(dist.Alpha+pos, dist.Beta+neg)
}

// Return the Beta which maximizes the probability of emitting the given sequence,
// based on a method of moments estimation
func (dist Beta) MaximizeByMoM(vals []float64) *Beta {
	var mean, variance float64
	for _, v := range vals {
		mean += v
	}
	mean /= float64(len(vals))
	for _, v := range vals {
		diff := mean - v
		variance += diff * diff
	}
	variance /= float64(len(vals)) - 1

	alpha := mean * (((mean * (1 - mean)) / variance) - 1)
	beta := (1 - mean) * (((mean * (1 - mean)) / variance) - 1)

	return NewBetaDist(alpha, beta)
}
