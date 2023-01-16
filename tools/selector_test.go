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
				// Keep all backups within the last 2 days

				// Then keep one backup every 7 days (until we hit 14 days)
				"@2d/7d",

				// Then keep one backup every 14 days (until we hit 30 days)
				"@14d/30d",

				// Then keep one backup every 30 days (until we hit 365 days)
				"@30d/365d",
			})

			now := time.Now().UTC()

			// We should always match the current time
			So(Matches(now, specs), ShouldBeTrue)

			// And we should always match yesterday at the same time
			So(Matches(now.Add(-24*time.Hour), specs), ShouldBeTrue)

			startOfDay := now.Truncate(24 * time.Hour)
			// We should match the start of the day
			So(Matches(startOfDay, specs), ShouldBeTrue)

			startOfYesterday := startOfDay.Add(-24 * time.Hour)
			// We should match the start of yesterday
			So(Matches(startOfYesterday, specs), ShouldBeTrue)

			// But we shouldn't match three days ago
			So(Matches(startOfYesterday.Add(-24*time.Hour).Add(-1*time.Second), specs), ShouldBeFalse)

			// We also want to keep a backup every 7 days if they're older than 2 days, so we should match the start of the week
			startOfWeek := now.Truncate(7 * 24 * time.Hour)
			So(Matches(startOfWeek, specs), ShouldBeTrue)

			// We should also consider cases where we're at the start of the current week (and include the prior week)
			startOfLastWeek := startOfWeek.Add(-7 * 24 * time.Hour)
			So(Matches(startOfLastWeek, specs), ShouldBeTrue)

			// But we shouldn't match an arbitrary backup from 3 weeks ago
			So(Matches(now.Add(-21*24*time.Hour), specs), ShouldBeFalse)
		})
	})
}
