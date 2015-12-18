package bsc

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/channel"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
	"math/rand"
)

// Create a new binary symmetric channel with the specified noise rates.
func NewBSCPair(noiseRate1, noiseRate2 float64) *BSCPair {
	if noiseRate1 < 0 || noiseRate1 > 1 {
		panic(stats.ErrfInvalidProb(noiseRate1))
	} else if noiseRate2 < 0 || noiseRate2 > 1 {
		panic(stats.ErrfInvalidProb(noiseRate2))
	}
	ch := &BSCPair{
		NoiseRate1: variable.NewContinuousRV(noiseRate1, dist.UnitIntervalSpace),
		NoiseRate2: variable.NewContinuousRV(noiseRate2, dist.UnitIntervalSpace),
	}
	ch.DefChannelSampleN.Channel = ch
	return ch
}

// A binary symmetric channel. Given a Bernoulli random variable X, it emits
// a Bernoulli random variable Y such that with probability `NoiseRate`
// Y = !X, and Y = X otherwise. In other words, the channel has fixed
// probability `NoiseRate` of outputting Y with the opposite value of X.
type BSCPair struct {

	// The probability of flipping the input for each layer of the channel
	NoiseRate1, NoiseRate2 *variable.ContinuousRV

	channel.DefChannelSampleN
}

// Send an input to the channel and sample an output
func (ch BSCPair) Sample(input variable.RandomVariable) variable.RandomVariable {
	var (
		rv    = input.(*variable.DiscreteRV)
		space = dist.BooleanSpace
		x     = space.BoolValue(rv.Outcome())
		flip1 = rand.Float64() <= ch.NoiseRate1.Val()
		flip2 = rand.Float64() <= ch.NoiseRate2.Val()
	)
	if flip1 != flip2 {
		return variable.NewDiscreteRV(space.BoolOutcome(!x), space)
	} else {
		return variable.NewDiscreteRV(rv.Outcome(), space)
	}
}

// Build a factor relating an input variable to an output variable
func (ch BSCPair) Factor(input variable.RandomVariable, output variable.RandomVariable) factor.Factor {
	return &BSCPairFactor{
		Input:      input.(*variable.DiscreteRV),
		Output:     output.(*variable.DiscreteRV),
		NoiseRate1: ch.NoiseRate1,
		NoiseRate2: ch.NoiseRate2,
	}
}

// Build factors relating an input variable to a sequence of output variables
func (ch BSCPair) Factors(input variable.RandomVariable, outputs []variable.RandomVariable) []factor.Factor {
	var fs []factor.Factor
	for _, rv := range outputs {
		fs = append(fs, ch.Factor(input, rv))
	}
	return fs
}

// A factor connecting an input variable to its output, as perturbed by a constant
// Bernoulli noise rate.
type BSCPairFactor struct {
	Input, Output          *variable.DiscreteRV
	NoiseRate1, NoiseRate2 *variable.ContinuousRV
}

// Do the input and output currently match?
func (factor BSCPairFactor) OutputMatchesInput() bool {
	return factor.Input.Equals(factor.Output)
}

// The adjacent random variables
func (factor BSCPairFactor) Adjacent() []variable.RandomVariable {
	return []variable.RandomVariable{factor.Output, factor.Input,
		factor.NoiseRate1, factor.NoiseRate2}
}

// The factor's current score, based on the values of adjacent variables
func (factor BSCPairFactor) Score() float64 {
	var n1, n2 = factor.NoiseRate1.Val(), factor.NoiseRate2.Val()
	if factor.OutputMatchesInput() {
		return (n1 * n2) + ((1 - n1) * (1 - n2))
	} else {
		return ((1 - n1) * n2) + (n1 * (1 - n2))
	}
}
