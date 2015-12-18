package dist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDenseMutableDiscreteDist(t *testing.T) {
	Convey("Test DenseMutableDiscreteDist interfaces", t, func() {
		dist := NewDenseMutableDiscreteDist(BooleanSpace)
		So(dist, ShouldImplement, (*MutableDiscreteDist)(nil))
	})
}
