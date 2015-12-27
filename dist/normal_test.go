package dist

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestNormal(t *testing.T) {
	Convey("Test Normal interfaces", t, func() {
		dist := NewNormalDist(1, 1)
		So(dist, ShouldImplement, (*Dist)(nil))
		So(dist, ShouldImplement, (*RealDist)(nil))
		So(dist, ShouldImplement, (*ContinuousDist)(nil))

		space := dist.Space()
		So(space, ShouldImplement, (*Space)(nil))
		So(space, ShouldImplement, (*RealSpace)(nil))
	})

	Convey("Test Normal dist", t, func() {
		dist := NewStandardNormalDist()
		So(dist.NumVars(), ShouldEqual, 1)
		So(dist.NumParams(), ShouldEqual, 2)
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.396952547477012)
		So(dist.Score([]float64{0.1}, []float64{0.5, 0.5}), ShouldAlmostEqual, 0.579383105522966)
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.396952547477012)
		dist.SetParams([]float64{0.5, 0.5})
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.579383105522966)
	})

	Convey("Test Normal PDF", t, func() {
		dist := NewNormalDist(0.1, 0.9)
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.443269200446036)
		So(dist.PDF(0.3), ShouldAlmostEqual, 0.432458299079718)
		So(dist.PDF(0.5), ShouldAlmostEqual, 0.401582033203048)
		So(dist.PDF(0.9), ShouldAlmostEqual, 0.298603179490360)

		dist = NewNormalDist(0.5, 0.5)
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.579383105522966)
		So(dist.PDF(0.3), ShouldAlmostEqual, 0.736540280606647)
		So(dist.PDF(0.5), ShouldAlmostEqual, 0.797884560802865)
		So(dist.PDF(0.9), ShouldAlmostEqual, 0.579383105522966)

		dist = NewNormalDist(0.9, 0.1)
		So(dist.PDF(0.1), ShouldAlmostEqual, 0.000000000000051)
		So(dist.PDF(0.3), ShouldAlmostEqual, 0.000000060758828)
		So(dist.PDF(0.5), ShouldAlmostEqual, 0.001338302257649)
		So(dist.PDF(0.9), ShouldAlmostEqual, 3.989422804014327)
	})

	Convey("Test Normal CDF", t, func() {
		dist := NewNormalDist(0.1, 0.9)
		So(dist.CDF(0.1), ShouldAlmostEqual, 0.5)
		So(dist.CDF(0.3), ShouldAlmostEqual, 0.587929552129057)
		So(dist.CDF(0.5), ShouldAlmostEqual, 0.671639356718115)
		So(dist.CDF(0.9), ShouldAlmostEqual, 0.812968601254559)

		dist = NewNormalDist(0.5, 0.5)
		So(dist.CDF(0.1), ShouldAlmostEqual, 0.211855398583397)
		So(dist.CDF(0.3), ShouldAlmostEqual, 0.344578258389676)
		So(dist.CDF(0.5), ShouldAlmostEqual, 0.5)
		So(dist.CDF(0.9), ShouldAlmostEqual, 0.788144601416603)

		dist = NewNormalDist(0.9, 0.1)
		So(dist.CDF(0.1), ShouldAlmostEqual, 0.000000000000001)
		So(dist.CDF(0.3), ShouldAlmostEqual, 0.000000000986588)
		So(dist.CDF(0.5), ShouldAlmostEqual, 0.000031671241833)
		So(dist.CDF(0.9), ShouldAlmostEqual, 0.5)
	})

	Convey("Test Normal mean", t, func() {
		So(NewNormalDist(0.1, 0.9).Mean(), ShouldAlmostEqual, 0.1)
		So(NewNormalDist(0.5, 0.5).Mean(), ShouldAlmostEqual, 0.5)
		So(NewNormalDist(0.9, 0.1).Mean(), ShouldAlmostEqual, 0.9)
	})

	Convey("Test Normal variance", t, func() {
		So(NewNormalDist(0.1, 0.9).Variance(), ShouldAlmostEqual, 0.81)
		So(NewNormalDist(0.5, 0.5).Variance(), ShouldAlmostEqual, 0.25)
		So(NewNormalDist(0.9, 0.1).Variance(), ShouldAlmostEqual, 0.01)
	})

	Convey("Test Normal draws", t, func() {
		const n = 100
		dist := NewNormalDist(1, 9)
		mean := 0.0
		for _, v := range dist.SampleN(n) {
			mean += v
		}
		mean /= n

		std := math.Sqrt(dist.Variance())
		So(mean, ShouldBeBetween, dist.Mean()-std, dist.Mean()+std)
	})
}
