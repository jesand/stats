package variable

import (
	"github.com/jesand/stats/dist"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestContinuousRV(t *testing.T) {
	Convey("Test ContinuousRV interfaces", t, func() {
		So(ContinuousRV{}, ShouldImplement, (*RandomVariable)(nil))
	})

	Convey("Test ContinuousRV", t, func() {
		rv := NewContinuousRV(0.5, dist.NewUnitIntervalSpace())
		So(rv.Val(), ShouldEqual, 0.5)
		So(rv.Space().Inf(), ShouldEqual, 0)
		So(rv.Space().Sup(), ShouldEqual, 1)
		rv.Set(0.4)
		So(rv.Val(), ShouldEqual, 0.4)
	})
}

func TestDiscreteRV(t *testing.T) {
	Convey("Test DiscreteRV interfaces", t, func() {
		So(DiscreteRV{}, ShouldImplement, (*RandomVariable)(nil))
	})

	Convey("Test DiscreteRV", t, func() {
		rv := NewDiscreteRV(1, dist.BooleanSpace{})
		So(rv.Val(), ShouldEqual, 1)
		So(rv.Space().Inf(), ShouldEqual, 0)
		So(rv.Space().Sup(), ShouldEqual, 1)
		rv.Set(0)
		So(rv.Val(), ShouldEqual, 0)
	})
}
