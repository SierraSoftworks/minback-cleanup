package tools

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSampleSpec(t *testing.T) {
	Convey("SampleSpec", t, func() {
		Convey("parseLongDuration", func() {

			Convey("1m", func() {
				d, err := parseLongDuration("1m")
				So(err, ShouldBeNil)
				So(d, ShouldEqual, time.Minute)
			})

			Convey("120m", func() {
				d, err := parseLongDuration("120m")
				So(err, ShouldBeNil)
				So(d, ShouldEqual, 120*time.Minute)
			})

			Convey("1h", func() {
				d, err := parseLongDuration("1h")
				So(err, ShouldBeNil)
				So(d, ShouldEqual, time.Hour)
			})

			Convey("5h", func() {
				d, err := parseLongDuration("5h")
				So(err, ShouldBeNil)
				So(d, ShouldEqual, 5*time.Hour)
			})

			Convey("1d", func() {
				d, err := parseLongDuration("1d")
				So(err, ShouldBeNil)
				So(d, ShouldEqual, 24*time.Hour)
			})

			Convey("7d", func() {
				d, err := parseLongDuration("7d")
				So(err, ShouldBeNil)
				So(d, ShouldEqual, 7*24*time.Hour)
			})

			Convey("1w2d", func() {
				d, err := parseLongDuration("1w2d")
				So(err, ShouldBeNil)
				So(d, ShouldEqual, 9*24*time.Hour)
			})
		})

		Convey("fmtLongDuration", func() {

			Convey("1m", func() {
				d := fmtLongDuration(time.Minute)
				So(d, ShouldEqual, "1m")
			})

			Convey("12m", func() {
				d := fmtLongDuration(12 * time.Minute)
				So(d, ShouldEqual, "12m")
			})

			Convey("1h", func() {
				d := fmtLongDuration(time.Hour)
				So(d, ShouldEqual, "1h")
			})

			Convey("5h", func() {
				d := fmtLongDuration(5 * time.Hour)
				So(d, ShouldEqual, "5h")
			})

			Convey("1d", func() {
				d := fmtLongDuration(24 * time.Hour)
				So(d, ShouldEqual, "1d")
			})

			Convey("5d", func() {
				d := fmtLongDuration(5 * 24 * time.Hour)
				So(d, ShouldEqual, "5d")
			})

			Convey("1w2d", func() {
				d := fmtLongDuration(9 * 24 * time.Hour)
				So(d, ShouldEqual, "1w2d")
			})
		})

		Convey("Serialization", func() {
			examples := []struct {
				Spec  string
				Value SampleSpec
			}{
				{"@1d/5d", SampleSpec{Before: 24 * time.Hour, Interval: 5 * 24 * time.Hour}},
				{"@1d/1w", SampleSpec{Before: 24 * time.Hour, Interval: 7 * 24 * time.Hour}},
				{"@1d/1w~1d", SampleSpec{Before: 24 * time.Hour, Interval: 7 * 24 * time.Hour, Round: 24 * time.Hour}},
				{"/1d", SampleSpec{Interval: 24 * time.Hour}},
			}

			Convey("MarshalString", func() {
				for _, example := range examples {
					Convey(example.Spec, func() {
						str, err := example.Value.MarshalString()
						So(err, ShouldBeNil)
						So(str, ShouldEqual, example.Spec)
					})
				}
			})

			Convey("UnmarshalString", func() {
				for _, example := range examples {
					Convey(example.Spec, func() {
						s := SampleSpec{}
						So(s.UnmarshalString(example.Spec), ShouldBeNil)
						So(s, ShouldResemble, example.Value)
					})
				}
			})

			Convey("MarshalJSON", func() {
				s := SampleSpec{
					Before:   24 * time.Hour,
					Interval: 7 * 24 * time.Hour,
				}

				b, err := s.MarshalJSON()
				So(err, ShouldBeNil)
				So(string(b), ShouldEqual, `"@1d/1w"`)
			})

			Convey("UnmarshalJSON", func() {
				s := SampleSpec{}

				So(s.UnmarshalJSON([]byte(`"@1d/7d"`)), ShouldBeNil)
				So(s.Before, ShouldEqual, 24*time.Hour)
				So(s.Interval, ShouldEqual, 7*24*time.Hour)
			})

		})

		Convey("WithinWindow", func() {
			Convey("5m~1m", func() {
				s, err := NewSampleSpec("/5m~1m")
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)

				So(s.WithinWindow(time.Now().Round(5*time.Minute).Add(-time.Minute)), ShouldBeTrue)
				So(s.WithinWindow(time.Now().Round(5*time.Minute).Add(-5*time.Second)), ShouldBeTrue)
			})

			Convey("@1h/5m", func() {
				s, err := NewSampleSpec("@1h/5m")
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)

				So(s.WithinWindow(time.Now().Round(time.Hour)), ShouldBeFalse)
				So(s.WithinWindow(time.Now().Add(-50*time.Minute)), ShouldBeFalse)
				So(s.WithinWindow(time.Now().Add(-2*time.Hour)), ShouldBeTrue)
			})

			Convey("@8h/15m", func() {
				s, err := NewSampleSpec("@8h/15m")
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)

				So(s.WithinWindow(time.Now().Round(7*time.Hour)), ShouldBeFalse)
				So(s.WithinWindow(time.Now().Add(-6*time.Hour)), ShouldBeFalse)
				So(s.WithinWindow(time.Now().Add(-9*time.Hour)), ShouldBeTrue)
			})
		})

		Convey("Matches", func() {
			Convey("5m~1m", func() {
				s, err := NewSampleSpec("/5m~1m")
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)

				So(s.Matches(time.Now().Round(5*time.Minute).Add(-time.Minute)), ShouldBeFalse)
				So(s.Matches(time.Now().Round(5*time.Minute).Add(5*time.Second)), ShouldBeTrue)
			})

			Convey("@1h/5m", func() {
				s, err := NewSampleSpec("@1h/5m")
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)

				So(s.Matches(time.Now().Round(time.Hour)), ShouldBeFalse)
				So(s.Matches(time.Now().Add(-2*time.Hour).Round(time.Hour)), ShouldBeTrue)
				So(s.Matches(time.Now().Add(-2*time.Hour).Round(5*time.Minute)), ShouldBeTrue)
				So(s.Matches(time.Now().Add(-2*time.Hour).Round(5*time.Minute).Add(time.Minute)), ShouldBeFalse)
			})

			Convey("@8h/15m", func() {
				s, err := NewSampleSpec("@8h/15m")
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)

				So(s.Matches(time.Now().Round(7*time.Hour)), ShouldBeFalse)
				So(s.Matches(time.Now().Add(-9*time.Hour).Round(time.Hour)), ShouldBeTrue)
				So(s.Matches(time.Now().Add(-30*time.Hour).Round(15*time.Minute)), ShouldBeTrue)
				So(s.Matches(time.Now().Add(-17*time.Hour).Round(15*time.Minute).Add(5*time.Minute)), ShouldBeFalse)
			})
		})
	})
}
