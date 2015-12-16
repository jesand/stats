package gibbs

import (
	// "github.com/jesand/stats/factor"
	// "github.com/jesand/stats/process"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGibbsSampling(t *testing.T) {
	SkipConvey("Given data from a Bernoulli process", t, func() {
		// var (
		// 	proc        = process.NewBernoulliProcess(0.7)
		// 	rvs         = proc.SampleN(10)
		// 	factors     = proc.Factors(rvs)
		// 	bernSampler = DistSampler{proc.Dist}
		// 	biasSampler = DistSampler{Beta}
		// 	model []GibbsSample
		// )
		// for i := range rvs {
		// 	model = append(model, GibbsSample{
		// 		Variable: rvs[i],
		// 		Factors:  []factor.Factor{factors[i]},
		// 		Sampler:  sampler,
		// 	})
		// }
		// vals := InferByGibbsSampling(model, 1000, 100)
		// for _, rv := range vals {

		// }
	})
}
