package process

import (
	"github.com/jesand/stats/variable"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBernoulliProcess(t *testing.T) {
	Convey("Test BernoulliProcess interfaces", t, func() {
		process := NewBernoulliProcess(0.5)
		So(process, ShouldImplement, (*StochasticProcess)(nil))
	})

	Convey("Test BernoulliProcess.Sample", t, func() {
		const n = 100
		process := NewBernoulliProcess(0.7)

		mean := 0.0
		for i := 0; i < n; i++ {
			rv := process.Sample()
			So(rv, ShouldHaveSameTypeAs, (*variable.DiscreteRV)(nil))
			mean += rv.Val()
		}
		mean /= n

		So(mean, ShouldBeBetween, 0.65, 0.75)
	})

	Convey("Test BernoulliProcess.SetBias", t, func() {
		const n = 100
		process := NewBernoulliProcess(0.7)
		process.SetBias(0.5)

		mean := 0.0
		for i := 0; i < n; i++ {
			rv := process.Sample()
			So(rv, ShouldHaveSameTypeAs, (*variable.DiscreteRV)(nil))
			mean += rv.Val()
		}
		mean /= n

		So(mean, ShouldBeBetween, 0.45, 0.55)
	})

	Convey("Test BernoulliProcess.SampleN", t, func() {
		const n = 100
		process := NewBernoulliProcess(0.7)

		mean := 0.0
		for _, rv := range process.SampleN(n) {
			So(rv, ShouldHaveSameTypeAs, (*variable.DiscreteRV)(nil))
			mean += rv.Val()
		}
		mean /= n

		So(mean, ShouldBeBetween, 0.65, 0.75)
	})

	Convey("Test BernoulliProcess.Factor", t, func() {
		process := NewBernoulliProcess(0.7)
		bias := process.Params[0]
		rvs := process.SampleN(10)
		factors := process.Factors(rvs)
		for i, factor := range factors {
			So(factor.Adjacent(), ShouldResemble, []variable.RandomVariable{rvs[i], bias})
			if rvs[i].Val() == 0 {
				So(factor.Score(), ShouldAlmostEqual, 0.3)
			} else {
				So(factor.Score(), ShouldAlmostEqual, 0.7)
			}
		}

		process.SetBias(0.3)
		So(bias.Val(), ShouldEqual, 0.3)
		for i, factor := range factors {
			if rvs[i].Val() == 0 {
				So(factor.Score(), ShouldAlmostEqual, 0.7)
			} else {
				So(factor.Score(), ShouldAlmostEqual, 0.3)
			}
		}

		bias.Set(0.7)
		for i, factor := range factors {
			if rvs[i].Val() == 0 {
				So(factor.Score(), ShouldAlmostEqual, 0.3)
			} else {
				So(factor.Score(), ShouldAlmostEqual, 0.7)
			}
		}
	})
}
