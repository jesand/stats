package model

import (
	"github.com/jesand/stats/channel/bsc"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
	"math"
)

// Create a new MultipleBSCPairModel
func NewMultipleBSCPairModel() *MultipleBSCPairModel {
	return &MultipleBSCPairModel{
		Noise1Rates: make(map[string]*variable.ContinuousRV),
		Noise2Rates: make(map[string]*variable.ContinuousRV),
		Channels:    make(map[string]map[string]*bsc.BSCPair),
		Inputs:      make(map[string]*variable.DiscreteRV),
		FactorGraph: factor.NewFactorGraph(),
		SoftInputs:  true,
	}
}

// A noisy channel model which explains a data set as passing through one of
// a set of BSC channels with unknown noise rates. For instance, this can be
// used to infer the answers of independent binary-valued crowdsourcing
// questions.
type MultipleBSCPairModel struct {

	// The noise rates
	Noise1Rates, Noise2Rates map[string]*variable.ContinuousRV

	// The variables sent over the noisy channels
	Inputs map[string]*variable.DiscreteRV

	// The posterior probability that each latent input variable is true
	InputScores map[string]float64

	// Whether to use soft or hard assignments to inputs during inference
	SoftInputs bool

	// The noisy channels
	Channels map[string]map[string]*bsc.BSCPair

	// A factor graph relating inputs to outputs
	FactorGraph *factor.FactorGraph
}

// Adds a new BSCPair to the model with the given layer names and noise rates.
// If a given noise rate has already been added, it will not be changed.
func (model *MultipleBSCPairModel) AddChannel(name1 string, noise1 float64,
	name2 string, noise2 float64) {

	ch := bsc.NewBSCPair(noise1, noise2)

	n1, ok := model.Noise1Rates[name1]
	if ok {
		ch.NoiseRate1 = n1
	} else {
		model.Noise1Rates[name1] = ch.NoiseRate1
	}

	n2, ok := model.Noise2Rates[name2]
	if ok {
		ch.NoiseRate2 = n2
	} else {
		model.Noise2Rates[name2] = ch.NoiseRate2
	}

	if _, ok := model.Channels[name1]; !ok {
		model.Channels[name1] = make(map[string]*bsc.BSCPair)
	}
	model.Channels[name1][name2] = ch
}

// Ask whether a given channel exists
func (model MultipleBSCPairModel) HasChannel(name1, name2 string) bool {
	chs, ok := model.Channels[name1]
	if !ok {
		return false
	}
	_, ok = chs[name2]
	return ok
}

// Ask whether a given input exists
func (model MultipleBSCPairModel) HasInput(name string) bool {
	_, ok := model.Inputs[name]
	return ok
}

// Adds a new observation to the model for the given channel and input. If the
// input is new, it will be created automatically.
func (model *MultipleBSCPairModel) AddObservation(input, channel1, channel2 string, value bool) {
	inputVar, ok := model.Inputs[input]
	if !ok {
		inputVar = variable.NewDiscreteRV(0, dist.BooleanSpace)
		model.Inputs[input] = inputVar
	}
	ch := model.Channels[channel1][channel2]

	// Add a worker noise factor to explain the assessment
	model.FactorGraph.AddFactor(ch.Factor(inputVar, variable.NewDiscreteRV(
		dist.BooleanSpace.BoolOutcome(value), dist.BooleanSpace)))
}

// Score the model, using the current parameter values
func (model MultipleBSCPairModel) Score() float64 {
	return model.FactorGraph.Score()
}

// Train noise rates and input values using expectation maximization.
func (model *MultipleBSCPairModel) EM(maxRounds int, tolerance float64,
	callback func(model *MultipleBSCPairModel, round int, stage string)) {

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

		thisRound2, lastRound2 := thisRound, lastRound
		for r2 := 1; (maxRounds == 0 || r2 <= maxRounds) &&
			thisRound2-lastRound2 > tolerance; r2++ {

			// Update the first layer of noise rates
			for _, noiseRate := range model.Noise1Rates {
				var count, sum float64
				for _, factor := range model.FactorGraph.AdjToVariable(noiseRate) {
					if ch, ok := factor.(*bsc.BSCPairFactor); ok {
						count++
						qi := softScores[ch.Input]
						n2 := ch.NoiseRate2.Val()
						if ch.Output.Val() == 1 {
							sum += qi*n2 + (1-qi)*(1-n2)
						} else {
							sum += qi*(1-n2) + (1-qi)*n2
						}
					}
				}
				if sum == 0 {
					noiseRate.Set(1e-3)
				} else if sum == count {
					noiseRate.Set(1 - 1e-3)
				} else {
					noiseRate.Set(sum / count)
				}
			}
			if callback != nil {
				callback(model, round, "noise1")
			}

			// Update the second layer of noise rates
			for _, noiseRate := range model.Noise2Rates {
				var count, sum float64
				for _, factor := range model.FactorGraph.AdjToVariable(noiseRate) {
					if ch, ok := factor.(*bsc.BSCPairFactor); ok {
						count++
						qi := softScores[ch.Input]
						n1 := ch.NoiseRate1.Val()
						if ch.Output.Val() == 1 {
							sum += qi*n1 + (1-qi)*(1-n1)
						} else {
							sum += qi*(1-n1) + (1-qi)*n1
						}
					}
				}
				if sum == 0 {
					noiseRate.Set(1e-3)
				} else if sum == count {
					noiseRate.Set(1 - 1e-3)
				} else {
					noiseRate.Set(sum / count)
				}
			}
			if callback != nil {
				callback(model, round, "noise2")
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
