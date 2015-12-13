package factor

import (
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/variable"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestDistFactor(t *testing.T) {
	Convey("Test DistFactor interfaces", t, func() {
		So(DistFactor{}, ShouldImplement, (*Factor)(nil))
	})

	Convey("Test DistFactor", t, func() {
		val := variable.NewContinuousRV(0.1, dist.NewUnitIntervalSpace())
		alpha := variable.NewContinuousRV(0.5, dist.NewRealIntervalSpace(0, math.Inf(+1)))
		beta := variable.NewContinuousRV(0.5, dist.NewRealIntervalSpace(0, math.Inf(+1)))
		factor := NewDistFactor([]variable.RandomVariable{val, alpha, beta}, dist.NewBetaDist(0, 0))
		So(factor.Adjacent(), ShouldResemble, []variable.RandomVariable{val, alpha, beta})
		So(factor.Score(), ShouldAlmostEqual, 0.0854694647)
	})
}
