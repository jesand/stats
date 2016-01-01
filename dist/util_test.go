package dist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtil(t *testing.T) {
	var (
		x1 = []float64{}
		x2 = []float64{1, 1, 1, 1}
		x3 = []float64{1, 2, 3, 4}
	)
	Convey("Test Mean()", t, func() {
		So(Mean(x1), ShouldEqual, 0)
		So(Mean(x2), ShouldEqual, 1)
		So(Mean(x3), ShouldEqual, 2.5)
	})
	Convey("Test Variance()", t, func() {
		So(Variance(x1), ShouldEqual, 0)
		So(Variance(x2), ShouldEqual, 0)
		So(Variance(x3), ShouldAlmostEqual, 1.666666666666667)
	})
	Convey("Test Min()", t, func() {
		So(Min(x1), ShouldEqual, 0)
		So(Min(x2), ShouldEqual, 1)
		So(Min(x3), ShouldEqual, 1)
	})
	Convey("Test MinGt()", t, func() {
		So(MinGt(x1, 0), ShouldEqual, 0)
		So(MinGt(x2, 0), ShouldEqual, 1)
		So(MinGt(x2, 1), ShouldEqual, 1)
		So(MinGt(x3, 1), ShouldEqual, 2)
	})
	Convey("Test Max()", t, func() {
		So(Max(x1), ShouldEqual, 0)
		So(Max(x2), ShouldEqual, 1)
		So(Max(x3), ShouldEqual, 4)
	})
	Convey("Test MaxLt()", t, func() {
		So(MaxLt(x1, 0), ShouldEqual, 0)
		So(MaxLt(x2, 0), ShouldEqual, 0)
		So(MaxLt(x2, 1), ShouldEqual, 1)
		So(MaxLt(x2, 2), ShouldEqual, 1)
		So(MaxLt(x3, 3), ShouldEqual, 2)
	})
}
