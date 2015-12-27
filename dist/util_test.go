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
}
