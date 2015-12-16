package gibbs

import (
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
)

// Indicates that a particular variable should be sampled with a particular
// sampler in the order provided to InferByGibbsSampling().
type GibbsSample struct {
	Variable variable.RandomVariable
	Factors  []factor.Factor
	Sampler  ValueSampler
}

// Select values for all latent variables using Gibbs sampling. We iterate
// over the model in the provided order. We run `burnin` iterations to allow
// the model to become calibrated, and then sample each variable in turn with
// `thinning` full rounds of sampling in between each variable's draw.
// Returns the sampled values for all variables, in the same order as specified
// in `model.`
func Infer(model []GibbsSample, burnin, thinning int) []variable.RandomVariable {

	// Burn-in period
	for r := 0; r < burnin; r++ {
		gibbsRound(model, -1)
	}

	// Draw variable values
	var state []variable.RandomVariable
	for vIdx := 0; vIdx < len(model); vIdx++ {
		for r := 0; r < thinning; r++ {
			gibbsRound(model, -1)
		}
		state = append(state, gibbsRound(model, vIdx))
	}
	return state
}

// Runs one full round of Gibbs sampling. If vIdx >= 0, then the value sampled
// for the variable with index vIdx will be returned. Otherwise returns nil.
func gibbsRound(model []GibbsSample, vIdx int) variable.RandomVariable {
	var output variable.RandomVariable
	for i, v := range model {
		v.Sampler.SampleValue(v.Variable, v.Factors)
		if i == vIdx {
			if cv, ok := v.Variable.(*variable.ContinuousRV); ok {
				output = variable.NewContinuousRV(cv.Val(), cv.Space())
			} else if dv, ok := v.Variable.(*variable.DiscreteRV); ok {
				output = variable.NewDiscreteRV(dv.Outcome(), dv.Space())
			}
		}
	}
	return output
}
