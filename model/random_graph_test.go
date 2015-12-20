package model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_RandomBipartiteGraph(t *testing.T) {
	Convey("RandomBipartiteGraph() handles invalid input", t, func() {
		var err error
		_, err = RandomBipartiteGraph(nil, nil)
		So(err.Error(), ShouldEqual, "Total node degree is zero")
		_, err = RandomBipartiteGraph([]int{1, 2, 3}, []int{3, 2, 2})
		So(err.Error(), ShouldEqual, "Total left degree 6 != total right degree 7")
	})
	Convey("RandomBipartiteGraph() returns a correct graph", t, func() {
		var (
			leftDegrees  = []int{2, 4, 6}
			rightDegrees = []int{2, 2, 2, 2, 2, 2}
		)
		edges, err := RandomBipartiteGraph(leftDegrees, rightDegrees)
		So(err, ShouldBeNil)
		So(len(edges), ShouldEqual, 12)
		for _, edge := range edges {
			So(edge[0], ShouldBeBetweenOrEqual, 0, len(leftDegrees)-1)
			So(edge[1], ShouldBeBetweenOrEqual, 0, len(rightDegrees)-1)
			leftDegrees[edge[0]]--
			rightDegrees[edge[1]]--
		}
		So(leftDegrees, ShouldResemble, []int{0, 0, 0})
		So(rightDegrees, ShouldResemble, []int{0, 0, 0, 0, 0, 0})
	})
}
