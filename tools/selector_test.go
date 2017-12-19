package tools

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func testBuildSpecs(c C, specs []string) []SampleSpec {
	out := []SampleSpec{}
	for _, spec := range specs {
		s, err := NewSampleSpec(spec)
		c.So(err, ShouldBeNil)
		out = append(out, s)
	}

	return out
}

func TestSelector(t *testing.T) {
	Convey("Selector", t, func() {
		Convey("Matches", func(c C) {
			specs := testBuildSpecs(c, []string{
				"~2d/7d",
				"~14d/30d",
				"~30d/365d",
			})

			So(Matches(time.Now().Round(24*time.Hour), specs), ShouldBeTrue)
			So(Matches(time.Now().Round(7*24*time.Hour).Add(-24*time.Hour), specs), ShouldBeFalse)
			So(Matches(time.Now().Round(7*24*time.Hour).Add(-7*24*time.Hour), specs), ShouldBeTrue)
		})
	})
}
