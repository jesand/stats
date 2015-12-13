package dist

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestBeta(t *testing.T) {
	Convey("Test Beta interfaces", t, func() {
		dist := NewBetaDist(1, 1)
		So(dist, ShouldImplement, (*Dist)(nil))
		So(dist, ShouldImplement, (*RealDist)(nil))
		So(dist, ShouldImplement, (*ContinuousDist)(nil))

		space := dist.Space()
		So(space, ShouldImplement, (*Space)(nil))
		So(space, ShouldImplement, (*RealSpace)(nil))
	})

	Convey("Test Beta dist", t, func() {
		dist := NewBetaDist(0.1, 0.9)
		So(dist.NumVars(), ShouldEqual, 1)
		So(dist.NumParams(), ShouldEqual, 2)
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.789602001365603)
		So(dist.Score([]float64{0.1}, []float64{0.5, 0.5}), ShouldAlmostEqual, 0.0854694647)
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.789602001365603)
		dist.SetParams([]float64{0.5, 0.5})
		So(dist.PDF(0.1), ShouldAlmostEqual, 1.061032953945969)
	})

	Convey("Test Beta PDF", t, func() {
		beta := NewBetaDist(0.1, 0.9)
		So(beta.PDF(0.1), ShouldAlmostEqual, 0.789602001365603)
		So(beta.PDF(0.3), ShouldAlmostEqual, 0.301240637595394)
		So(beta.PDF(0.5), ShouldAlmostEqual, 0.196726328616693)
		So(beta.PDF(0.9), ShouldAlmostEqual, 0.136148930108213)

		beta = NewBetaDist(0.5, 0.5)
		So(beta.PDF(0.1), ShouldAlmostEqual, 1.061032953945969)
		So(beta.PDF(0.3), ShouldAlmostEqual, 0.694609118042857)
		So(beta.PDF(0.5), ShouldAlmostEqual, 0.636619772367581)
		So(beta.PDF(0.9), ShouldAlmostEqual, 1.061032953945969)

		beta = NewBetaDist(0.9, 0.1)
		So(beta.PDF(0.1), ShouldAlmostEqual, 0.136148930108213)
		So(beta.PDF(0.3), ShouldAlmostEqual, 0.152943889294467)
		So(beta.PDF(0.5), ShouldAlmostEqual, 0.196726328616693)
		So(beta.PDF(0.9), ShouldAlmostEqual, 0.789602001365604)
	})

	Convey("Test Beta CDF", t, func() {
		beta := NewBetaDist(0.1, 0.9)
		So(beta.CDF(0.1), ShouldAlmostEqual, 0.782058177948575)
		So(beta.CDF(0.3), ShouldAlmostEqual, 0.874676058273383)
		So(beta.CDF(0.5), ShouldAlmostEqual, 0.922739224921369)
		So(beta.CDF(0.9), ShouldAlmostEqual, 0.985614974243009)

		beta = NewBetaDist(0.5, 0.5)
		So(beta.CDF(0.1), ShouldAlmostEqual, 0.204832764699133)
		So(beta.CDF(0.3), ShouldAlmostEqual, 0.369010119565545)
		So(beta.CDF(0.5), ShouldAlmostEqual, 0.500000000000000)
		So(beta.CDF(0.9), ShouldAlmostEqual, 0.795167235300867)

		beta = NewBetaDist(0.9, 0.1)
		So(beta.CDF(0.1), ShouldAlmostEqual, 0.014385025756991)
		So(beta.CDF(0.3), ShouldAlmostEqual, 0.042845131375421)
		So(beta.CDF(0.5), ShouldAlmostEqual, 0.077260775078631)
		So(beta.CDF(0.9), ShouldAlmostEqual, 0.217941822051425)
	})

	Convey("Test Beta mean", t, func() {
		So(NewBetaDist(0.1, 0.9).Mean(), ShouldAlmostEqual, 0.1)
		So(NewBetaDist(0.5, 0.5).Mean(), ShouldAlmostEqual, 0.5)
		So(NewBetaDist(0.9, 0.1).Mean(), ShouldAlmostEqual, 0.9)
	})

	Convey("Test Beta variance", t, func() {
		So(NewBetaDist(0.1, 0.9).Variance(), ShouldAlmostEqual, 0.045)
		So(NewBetaDist(0.5, 0.5).Variance(), ShouldAlmostEqual, 0.125)
		So(NewBetaDist(0.9, 0.1).Variance(), ShouldAlmostEqual, 0.045)
	})

	Convey("Test Beta draws", t, func() {
		const n = 100
		beta := NewBetaDist(1, 9)
		mean := 0.0
		for _, v := range beta.SampleN(n) {
			mean += v
		}
		mean /= n

		std := math.Sqrt(beta.Variance())
		So(mean, ShouldBeBetween, beta.Mean()-std, beta.Mean()+std)
	})
}
