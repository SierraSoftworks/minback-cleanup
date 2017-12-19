package tools

import (
	"testing"
	"time"

	"github.com/araddon/dateparse"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFilenameParser(t *testing.T) {
	Convey("ParseFilename", t, func() {
		examples := []struct {
			Filename string
			Prefix   string
			Time     time.Time
		}{
			{"2017-08-17.archive", "", dateparse.MustParse("2017-08-17")},
			{"postgres-2017-08-17.archive", "postgres-", dateparse.MustParse("2017-08-17")},
		}

		for _, example := range examples {
			Convey(example.Filename, func() {
				t, err := ParseFilename(example.Filename, example.Prefix)
				So(err, ShouldBeNil)
				So(t.Equal(example.Time), ShouldBeTrue)
			})
		}
	})
}
