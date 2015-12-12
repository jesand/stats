package dist

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestBernoulli(t *testing.T) {
	Convey("Test Bernoulli interfaces", t, func() {
		dist := NewBernoulliDist(0.5)
		So(dist, ShouldImplement, (*MutableDiscreteDist)(nil))
		So(dist, ShouldImplement, (*RealDist)(nil))
		So(dist.Space(), ShouldImplement, (*DiscreteRealSpace)(nil))
	})

	Convey("Test Bernoulli PMF", t, func() {
		dist := NewBernoulliDist(0.1)
		So(dist.Prob(0), ShouldAlmostEqual, 0.9)
		So(dist.Prob(1), ShouldAlmostEqual, 0.1)

		dist = NewBernoulliDist(0.5)
		So(dist.Prob(0), ShouldAlmostEqual, 0.5)
		So(dist.Prob(1), ShouldAlmostEqual, 0.5)

		dist = NewBernoulliDist(0.9)
		So(dist.Prob(0), ShouldAlmostEqual, 0.1)
		So(dist.Prob(1), ShouldAlmostEqual, 0.9)
	})

	Convey("Test Bernoulli CDF", t, func() {
		dist := NewBernoulliDist(0.1)
		So(dist.CDF(0), ShouldAlmostEqual, 0.9)
		So(dist.CDF(1), ShouldAlmostEqual, 1)

		dist = NewBernoulliDist(0.5)
		So(dist.CDF(0), ShouldAlmostEqual, 0.5)
		So(dist.CDF(1), ShouldAlmostEqual, 1)

		dist = NewBernoulliDist(0.9)
		So(dist.CDF(0), ShouldAlmostEqual, 0.1)
		So(dist.CDF(1), ShouldAlmostEqual, 1)
	})

	Convey("Test Bernoulli mean", t, func() {
		So(NewBernoulliDist(0.1).Mean(), ShouldAlmostEqual, 0.1)
		So(NewBernoulliDist(0.5).Mean(), ShouldAlmostEqual, 0.5)
		So(NewBernoulliDist(0.9).Mean(), ShouldAlmostEqual, 0.9)
	})

	Convey("Test Bernoulli variance", t, func() {
		So(NewBernoulliDist(0.1).Variance(), ShouldAlmostEqual, 0.09)
		So(NewBernoulliDist(0.5).Variance(), ShouldAlmostEqual, 0.25)
		So(NewBernoulliDist(0.9).Variance(), ShouldAlmostEqual, 0.09)
	})

	Convey("Test Bernoulli draws", t, func() {
		const n = 100
		dist := NewBernoulliDist(0.7)
		mean := 0.0
		for _, outcome := range dist.SampleN(n) {
			mean += dist.BSpace().F64Value(outcome)
		}
		mean /= n

		std := math.Sqrt(dist.Variance())
		So(mean, ShouldBeBetween, dist.Mean()-std, dist.Mean()+std)
	})
}
