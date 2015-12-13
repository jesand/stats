package process

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/variable"
)

// Generate a new BernoulliProcess
func NewBernoulliProcess(bias float64) *BernoulliProcess {
	if bias < 0 || bias > 1 {
		panic(stats.ErrfInvalidProb(bias))
	}
	biasRV := variable.NewContinuousRV(bias, dist.NewUnitIntervalSpace())
	return &BernoulliProcess{
		IIDProcess: NewIIDProcess(
			[]variable.RandomVariable{biasRV},
			dist.NewBernoulliDist(bias),
		),
	}
}

// A Bernoulli Process generates an infinite sequence of binary-valued random
// variables from the same Bernoulli distribution.
type BernoulliProcess struct {
	*IIDProcess
}

func (process *BernoulliProcess) SetBias(bias float64) {
	process.Params[0].Set(bias)
	process.Dist.(*dist.BernoulliDist).SetBias(bias)
}
