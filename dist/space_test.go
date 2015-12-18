package dist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnitInterval(t *testing.T) {
	var sp = UnitIntervalSpace
	Convey("Test UnitInterval interfaces", t, func() {
		So(sp, ShouldImplement, (*Space)(nil))
		So(sp, ShouldImplement, (*RealSpace)(nil))
	})
	Convey("Test UnitInterval.Inf()", t, func() {
		So(sp.Inf(), ShouldEqual, 0)
	})
	Convey("Test UnitInterval.Sup()", t, func() {
		So(sp.Sup(), ShouldEqual, 1)
	})
}

func TestBooleanSpace(t *testing.T) {
	var sp = BooleanSpace
	Convey("Test BooleanSpace interfaces", t, func() {
		So(sp, ShouldImplement, (*Space)(nil))
		So(sp, ShouldImplement, (*DiscreteSpace)(nil))
		So(sp, ShouldImplement, (*RealSpace)(nil))
		So(sp, ShouldImplement, (*DiscreteRealSpace)(nil))
	})
	Convey("Test BooleanSpace.Inf()", t, func() {
		So(sp.Inf(), ShouldEqual, 0)
	})
	Convey("Test BooleanSpace.Sup()", t, func() {
		So(sp.Sup(), ShouldEqual, 1)
	})
	Convey("Test BooleanSpace.Size()", t, func() {
		So(sp.Size(), ShouldEqual, 2)
	})
	Convey("Test BooleanSpace.F64Value()", t, func() {
		So(sp.F64Value(0), ShouldEqual, 0)
		So(sp.F64Value(1), ShouldEqual, 1)
		So(sp.F64Value(2), ShouldEqual, 1)
	})
	Convey("Test BooleanSpace.BoolValue()", t, func() {
		So(sp.BoolValue(0), ShouldEqual, false)
		So(sp.BoolValue(1), ShouldEqual, true)
		So(sp.BoolValue(2), ShouldEqual, true)
	})
}
