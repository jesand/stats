package model

import (
	"github.com/jesand/stats/channel/bsc"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
	"math"
)

// Create a new MultipleBSCModel
func NewMultipleBSCModel() *MultipleBSCModel {
	return &MultipleBSCModel{
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

// Adds a new BSC to the model with the given name and noise rate.
func (model *MultipleBSCModel) AddChannel(name string, noise float64) {
	model.Channels[name] = bsc.NewBSC(noise)
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
			ifFalse := math.Exp(model.FactorGraph.ScoreVar(input))
			input.Set(1)
			ifTrue := math.Exp(model.FactorGraph.ScoreVar(input))
			if ifFalse > ifTrue {
				input.Set(0)
			}
			if model.SoftInputs {
				if ifTrue == 0 {
					softScores[input] = 1e-6
				} else if ifFalse == 0 {
					softScores[input] = 1 - 1e-6
				} else {
					softScores[input] = ifTrue / (ifTrue + ifFalse)
				}
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
			}
			if callback != nil {
				callback(model, round, "noise")
			}

			lastRound2, thisRound2 = thisRound2, model.Score()
		}

		lastRound, thisRound = thisRound, model.Score()
	}

	model.InputScores = make(map[string]float64)
	for name, input := range model.Inputs {
		input.Set(0)
		ifFalse := math.Exp(model.FactorGraph.ScoreVar(input))
		input.Set(1)
		ifTrue := math.Exp(model.FactorGraph.ScoreVar(input))
		if ifFalse > ifTrue {
			input.Set(0)
		}
		model.InputScores[name] = ifTrue / (ifTrue + ifFalse)
	}
	if callback != nil {
		callback(model, 0, "Final")
	}
}
