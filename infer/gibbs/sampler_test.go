package gibbs

import (
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/process"
	"github.com/jesand/stats/variable"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProdValueSampler(t *testing.T) {
	Convey("Given a ProdValueSampler with a ConstFactor", t, func() {
		var (
			rv      = process.NewBernoulliProcess(0.5).Sample().(*variable.DiscreteRV)
			space   = dist.BooleanSpace
			factors = []factor.Factor{
				factor.NewConstFactor([]variable.RandomVariable{rv}, 20),
			}
			sampler = ProdValueSampler{}
		)
		Convey("The sampled value is uniformly random", func() {
			var pos, neg int
			for i := 0; i < 1000; i++ {
				sampler.SampleValue(rv, factors)
				if space.BoolValue(rv.Outcome()) {
					pos++
				} else {
					neg++
				}
			}
			So(pos, ShouldBeBetween, 400, 600)
			So(neg, ShouldBeBetween, 400, 600)
		})
	})
	Convey("Given a ProdValueSampler with 3 Bernoulli factors", t, func() {
		var (
			proc1   = process.NewBernoulliProcess(0.7)
			proc2   = process.NewBernoulliProcess(0.3)
			proc3   = process.NewBernoulliProcess(0.5)
			rv      = proc1.Sample().(*variable.DiscreteRV)
			rvs     = []variable.RandomVariable{rv}
			space   = dist.BooleanSpace
			factors = append(
				proc1.Factors(rvs),
				append(proc2.Factors(rvs),
					proc3.Factors(rvs)...)...)
			sampler = ProdValueSampler{}
		)
		Convey("The sampled value depends on the Bernoullis", func() {
			const margin = 25

			// Total prob: .7*.3*.5/(.7*.3*.5 + .3*.7*.5) = 1/2
			var pos, neg int
			for i := 0; i < 1000; i++ {
				sampler.SampleValue(rv, factors)
				if space.BoolValue(rv.Outcome()) {
					pos++
				} else {
					neg++
				}
			}
			So(pos, ShouldBeBetween, 500-margin, 500+margin)
			So(neg, ShouldBeBetween, 500-margin, 500+margin)

			// Total prob: .7*.6*.5/(.7*.6*.5 + .3*.4*.5) = 0.77
			proc2.SetBias(0.6)
			pos, neg = 0, 0
			for i := 0; i < 1000; i++ {
				sampler.SampleValue(rv, factors)
				if space.BoolValue(rv.Outcome()) {
					pos++
				} else {
					neg++
				}
			}
			So(pos, ShouldBeBetween, 770-margin, 770+margin)
			So(neg, ShouldBeBetween, 230-margin, 230+margin)

			// Total prob: .4*.3*.5/(.4*.3*.5 + .6*.7*.5) = 0.22
			proc1.SetBias(0.4)
			proc2.SetBias(0.3)
			pos, neg = 0, 0
			for i := 0; i < 1000; i++ {
				sampler.SampleValue(rv, factors)
				if space.BoolValue(rv.Outcome()) {
					pos++
				} else {
					neg++
				}
			}
			So(pos, ShouldBeBetween, 220-margin, 220+margin)
			So(neg, ShouldBeBetween, 780-margin, 780+margin)
		})
	})
}
