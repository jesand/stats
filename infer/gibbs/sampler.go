package gibbs

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
	"math/rand"
)

// A ValueSampler samples a new value for a random variable based on the scores
// of its adjacent factors. The variable's value is updated in the process.
type ValueSampler interface {
	SampleValue(v variable.RandomVariable, factors []factor.Factor)
}

// DistSampler is a ValueSampler which samples from some distribution.
type DistSampler struct {
	Dist dist.Dist
}

func (sampler DistSampler) SampleValue(v variable.RandomVariable, factors []factor.Factor) {
	if dd, ok := sampler.Dist.(dist.DiscreteDist); ok {
		v.(*variable.DiscreteRV).SetOutcome(dd.Sample())
	} else if cd, ok := sampler.Dist.(dist.ContinuousDist); ok {
		v.(*variable.ContinuousRV).Set(cd.Sample())
	} else {
		panic(stats.ErrfUnsupportedDist(sampler.Dist))
	}
}

// A ValueSampler which samples discrete values in proportion to the product of all factors.
type ProdValueSampler struct{}

func (sampler ProdValueSampler) SampleValue(v variable.RandomVariable, factors []factor.Factor) {
	dv, ok := v.(*variable.DiscreteRV)
	if !ok {
		panic(stats.ErrDiscreteOnly)
	}
	var (
		props = make([]float64, dv.Space().Size())
		total float64
	)
	for i := range props {
		dv.SetOutcome(dist.Outcome(i))
		props[i] = 1
		for _, f := range factors {
			props[i] *= f.Score()
		}
		total += props[i]
	}

	var remaining = rand.Float64() * total
	for i, prop := range props {
		remaining -= prop
		if remaining <= 0 {
			dv.SetOutcome(dist.Outcome(i))
			return
		}
	}
}
