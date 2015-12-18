/*
 * Implements variations on a binary symmetric channel (BSC) with constant noise
 * rate p, which takes a "true" boolean value X as input and outputs X with
 * probability (1-p) and outputs !X with probability p.
 *
 * We send the many X values across many such channels, and attempt to infer
 * their respective noise rates and the values of the input variables based
 * on the rate of agreement between the channels.
 */
package bsc

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/channel"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
	"math/rand"
)

// Create a new binary symmetric channel with the specified noise rate.
func NewBSC(noiseRate float64) *BSC {
	if noiseRate < 0 || noiseRate > 1 {
		panic(stats.ErrfInvalidProb(noiseRate))
	}
	ch := &BSC{
		NoiseRate: variable.NewContinuousRV(noiseRate, dist.UnitIntervalSpace),
	}
	ch.DefChannelSampleN.Channel = ch
	return ch
}

// A binary symmetric channel. Given a Bernoulli random variable X, it emits
// a Bernoulli random variable Y such that with probability `NoiseRate`
// Y = !X, and Y = X otherwise. In other words, the channel has fixed
// probability `NoiseRate` of outputting Y with the opposite value of X.
type BSC struct {

	// The probability of flipping the input
	NoiseRate *variable.ContinuousRV

	channel.DefChannelSampleN
}

// Send an input to the channel and sample an output
func (ch BSC) Sample(input variable.RandomVariable) variable.RandomVariable {
	var (
		rv    = input.(*variable.DiscreteRV)
		space = dist.BooleanSpace
		x     = space.BoolValue(rv.Outcome())
	)
	if rand.Float64() <= ch.NoiseRate.Val() {
		return variable.NewDiscreteRV(space.BoolOutcome(!x), space)
	} else {
		return variable.NewDiscreteRV(rv.Outcome(), space)
	}
}

// Build a factor relating an input variable to an output variable
func (ch BSC) Factor(input variable.RandomVariable, output variable.RandomVariable) factor.Factor {
	return &BSCFactor{
		Input:     input.(*variable.DiscreteRV),
		Output:    output.(*variable.DiscreteRV),
		NoiseRate: ch.NoiseRate,
	}
}

// Build factors relating an input variable to a sequence of output variables
func (ch BSC) Factors(input variable.RandomVariable, outputs []variable.RandomVariable) []factor.Factor {
	var fs []factor.Factor
	for _, rv := range outputs {
		fs = append(fs, ch.Factor(input, rv))
	}
	return fs
}

// A factor connecting an input variable to its output, as perturbed by a constant
// Bernoulli noise rate.
type BSCFactor struct {
	Input, Output *variable.DiscreteRV
	NoiseRate     *variable.ContinuousRV
}

// Do the input and output currently match?
func (factor BSCFactor) OutputMatchesInput() bool {
	return factor.Input.Equals(factor.Output)
}

// The adjacent random variables
func (factor BSCFactor) Adjacent() []variable.RandomVariable {
	return []variable.RandomVariable{factor.Output, factor.Input, factor.NoiseRate}
}

// The factor's current score, based on the values of adjacent variables
func (factor BSCFactor) Score() float64 {
	if factor.OutputMatchesInput() {
		return 1 - factor.NoiseRate.Val()
	} else {
		return factor.NoiseRate.Val()
	}
}
