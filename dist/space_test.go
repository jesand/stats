package dist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnitInterval(t *testing.T) {
	var dist = NewUnitIntervalSpace()
	Convey("Test UnitInterval interfaces", t, func() {
		So(dist, ShouldImplement, (*Space)(nil))
		So(dist, ShouldImplement, (*RealSpace)(nil))
	})
	Convey("Test UnitInterval.Inf()", t, func() {
		So(dist.Inf(), ShouldEqual, 0)
	})
	Convey("Test UnitInterval.Sup()", t, func() {
		So(dist.Sup(), ShouldEqual, 1)
	})
}

func TestBooleanSpace(t *testing.T) {
	var dist BooleanSpace
	Convey("Test BooleanSpace interfaces", t, func() {
		So(dist, ShouldImplement, (*Space)(nil))
		So(dist, ShouldImplement, (*DiscreteSpace)(nil))
		So(dist, ShouldImplement, (*RealSpace)(nil))
		So(dist, ShouldImplement, (*DiscreteRealSpace)(nil))
	})
	Convey("Test BooleanSpace.Inf()", t, func() {
		So(dist.Inf(), ShouldEqual, 0)
	})
	Convey("Test BooleanSpace.Sup()", t, func() {
		So(dist.Sup(), ShouldEqual, 1)
	})
	Convey("Test BooleanSpace.Size()", t, func() {
		So(dist.Size(), ShouldEqual, 2)
	})
	Convey("Test BooleanSpace.F64Value()", t, func() {
		So(dist.F64Value(0), ShouldEqual, 0)
		So(dist.F64Value(1), ShouldEqual, 1)
		So(dist.F64Value(2), ShouldEqual, 1)
	})
	Convey("Test BooleanSpace.BoolValue()", t, func() {
		So(dist.BoolValue(0), ShouldEqual, false)
		So(dist.BoolValue(1), ShouldEqual, true)
		So(dist.BoolValue(2), ShouldEqual, true)
	})
}
