package model

import (
	"github.com/jesand/stats/channel/bsc"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
)

// Create a new MultipleBSCModel with the given Beta prior
func NewMultipleBSCModel(alpha, beta float64) *MultipleBSCModel {
	return &MultipleBSCModel{
		NoiseAlpha:  variable.NewContinuousRV(alpha, dist.PositiveRealSpace),
		NoiseBeta:   variable.NewContinuousRV(beta, dist.PositiveRealSpace),
		NoiseDist:   dist.NewBetaDist(alpha, beta),
		UpdateBeta:  false,
		Inputs:      make(map[string]*variable.DiscreteRV),
		Channels:    make(map[string]*bsc.BSC),
		FactorGraph: factor.NewFactorGraph(),
		SoftInputs:  true,
	}
}

// A noisy channel model which explains a data set as passing through one of
// a set of BSC channels with unknown noise rates. For instance, this can be
// used to infer the answers of independent binary-valued crowdsourcing
// questions.
type MultipleBSCModel struct {

	// The parameters on a Beta prior for channel noise parameters
	NoiseAlpha, NoiseBeta *variable.ContinuousRV

	// The Beta prior for channel noise
	NoiseDist *dist.Beta

	// Indicates whether the Beta parameters should be updated after each round
	UpdateBeta bool

	// The variables sent over the noisy channels
	Inputs map[string]*variable.DiscreteRV

	// The posterior probability that each latent input variable is true
	InputScores map[string]float64

	// Whether to use soft or hard assignments to inputs during inference
	SoftInputs bool

	// The noisy channels
	Channels map[string]*bsc.BSC

	// A factor graph relating inputs to outputs
	FactorGraph *factor.FactorGraph
}

// Adds a new BSC to the model with the given name and noise rate. If noise
// is zero, we sample from the beta prior.
func (model *MultipleBSCModel) AddChannel(name string, noise float64) {
	ch := bsc.NewBSC(model.NoiseDist.Sample())
	if noise != 0 {
		ch.NoiseRate.Set(noise)
	}
	model.Channels[name] = ch
	model.FactorGraph.AddFactor(factor.NewDistFactor(
		[]variable.RandomVariable{ch.NoiseRate, model.NoiseAlpha,
			model.NoiseBeta}, model.NoiseDist))
}

// Ask whether a given channel exists
func (model MultipleBSCModel) HasChannel(name string) bool {
	_, ok := model.Channels[name]
	return ok
}

// Ask whether a given input exists
func (model MultipleBSCModel) HasInput(name string) bool {
	_, ok := model.Inputs[name]
	return ok
}

// Adds a new observation to the model for the given channel and input. If the
// input is new, it will be created automatically.
func (model *MultipleBSCModel) AddObservation(input, channel string, value bool) {
	inputVar, ok := model.Inputs[input]
	if !ok {
		inputVar = variable.NewDiscreteRV(0, dist.BooleanSpace)
		model.Inputs[input] = inputVar
	}
	ch := model.Channels[channel]

	// Add a worker noise factor to explain the assessment
	outcome := dist.BooleanSpace.BoolOutcome(value)
	model.FactorGraph.AddFactor(ch.Factor(inputVar, variable.NewDiscreteRV(
		outcome, dist.BooleanSpace)))
}

// Score the model, using the current parameter values
func (model MultipleBSCModel) Score() float64 {
	return model.FactorGraph.Score()
}

// Train noise rates and input values using expectation maximization.
func (model *MultipleBSCModel) EM(maxRounds int, tolerance float64,
	callback func(model *MultipleBSCModel, round int, stage string)) {

	var (
		round        int
		initialScore = model.Score()
		thisRound    = initialScore
		lastRound    = thisRound - 1.0
		softScores   = make(map[*variable.DiscreteRV]float64)
	)

	if callback != nil {
		callback(model, round, "Initial")
	}
	for round = 1; (maxRounds == 0 || round <= maxRounds) &&
		thisRound-lastRound > tolerance; round++ {

		// Update input
		for _, input := range model.Inputs {
			input.Set(0)
			ifFalse := model.FactorGraph.ScoreVar(input)
			input.Set(1)
			ifTrue := model.FactorGraph.ScoreVar(input)
			if ifFalse > ifTrue {
				input.Set(0)
			}
			if model.SoftInputs {
				softScores[input] = ifTrue / (ifTrue + ifFalse)
			} else {
				softScores[input] = input.Val()
			}
		}
		if callback != nil {
			callback(model, round, "input")
		}

		// Update noise rates
		thisRound2, lastRound2 := thisRound, lastRound
		for r2 := 1; (maxRounds == 0 || r2 <= maxRounds) &&
			thisRound2-lastRound2 > tolerance; r2++ {
			var rates []float64
			for _, ch := range model.Channels {
				var sum, count float64
				for _, factor := range model.FactorGraph.AdjToVariable(ch.NoiseRate) {
					if ch, ok := factor.(*bsc.BSCFactor); ok {
						count++
						qi := softScores[ch.Input]
						if ch.Output.Val() == 1 {
							sum += qi
						} else {
							sum += 1 - qi
						}
					}
				}
				if sum == 0 {
					ch.NoiseRate.Set(1e-3)
				} else if sum == count {
					ch.NoiseRate.Set(1 - 1e-3)
				} else {
					ch.NoiseRate.Set(sum / count)
				}
				// ch.NoiseRate.Set((failures + model.NoiseAlpha.Val()) /
				// 	(tries + model.NoiseAlpha.Val() + model.NoiseBeta.Val()))
				rates = append(rates, ch.NoiseRate.Val())
			}
			if callback != nil {
				callback(model, round, "noise")
			}

			// Update Beta prior
			if model.UpdateBeta {
				model.NoiseDist = model.NoiseDist.MaximizeByMoM(rates)
				model.NoiseAlpha.Set(model.NoiseDist.Alpha)
				model.NoiseBeta.Set(model.NoiseDist.Beta)
				if callback != nil {
					callback(model, round, "beta")
				}
			}
			lastRound2, thisRound2 = thisRound2, model.Score()
		}

		lastRound, thisRound = thisRound, model.Score()
	}

	model.InputScores = make(map[string]float64)
	for name, input := range model.Inputs {
		input.Set(0)
		ifFalse := model.FactorGraph.ScoreVar(input)
		input.Set(1)
		ifTrue := model.FactorGraph.ScoreVar(input)
		if ifFalse > ifTrue {
			input.Set(0)
		}
		model.InputScores[name] = ifTrue / (ifTrue + ifFalse)
	}
	if callback != nil {
		callback(model, 0, "Final")
	}
}
