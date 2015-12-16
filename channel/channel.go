package channel

import (
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
)

// A generic channel: given a message X, it emits a random variable Y derived
// from X based on the properties of the channel.
type Channel interface {

	// Send an input to the channel and sample an output
	Sample(input variable.RandomVariable) variable.RandomVariable

	// Send a sequence of inputs and sample an output for each
	SampleN(inputs []variable.RandomVariable) []variable.RandomVariable

	// Build factors relating an input variable to a sequence of output variables
	Factors(input variable.RandomVariable, outputs []variable.RandomVariable) []factor.Factor
}

// A default implementation of SampleN() for a Channel
type DefChannelSampleN struct{ Channel Channel }

func (def DefChannelSampleN) SampleN(inputs []variable.RandomVariable) []variable.RandomVariable {
	var outputs []variable.RandomVariable
	for _, rv := range inputs {
		outputs = append(outputs, def.Channel.Sample(rv))
	}
	return outputs
}
