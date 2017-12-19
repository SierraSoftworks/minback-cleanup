package tools

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SampleSpec describes how backups should be sampled
// for a given time range
type SampleSpec struct {
	Before   time.Duration
	Interval time.Duration
}

func NewSampleSpec(spec string) (SampleSpec, error) {
	s := SampleSpec{}
	if err := s.UnmarshalString(spec); err != nil {
		return s, err
	}

	return s, nil
}

// WithinWindow determines whether a given time object
// is within the sampling window for this SamplingSpec
func (s *SampleSpec) WithinWindow(t time.Time) bool {
	absBefore := time.Now().Add(-s.Before)
	if t.After(absBefore) {
		return false
	}

	return true
}

// Matches determines whether a given time object
// matches a sampling specification
func (s *SampleSpec) Matches(t time.Time) bool {
	if !s.WithinWindow(t) {
		return false
	}

	if !t.Round(s.Interval).Equal(t) {
		return false
	}

	return true
}

func (s *SampleSpec) MarshalJSON() ([]byte, error) {
	str, err := s.MarshalString()
	if err != nil {
		return nil, err
	}

	return json.Marshal(str)
}

func (s *SampleSpec) MarshalString() (string, error) {
	return fmt.Sprintf("~%s/%s", fmtLongDuration(s.Before), fmtLongDuration(s.Interval)), nil
}

func (s *SampleSpec) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	return s.UnmarshalString(str)
}

func (s *SampleSpec) UnmarshalString(str string) error {
	if str[0] != '~' {
		return fmt.Errorf("tools: invalid sample spec %s", str)
	}

	splitIndex := strings.IndexByte(str, '/')
	if splitIndex == -1 {
		return fmt.Errorf("tools: invalid sample spec %s", str)
	}

	before, err := parseLongDuration(str[1:splitIndex])
	if err != nil {
		return fmt.Errorf("tools: invalid sample spec %s", str)
	}

	interval, err := parseLongDuration(str[splitIndex+1:])
	if err != nil {
		return fmt.Errorf("tools: invalid sample spec %s", str)
	}

	s.Before = before
	s.Interval = interval

	return nil
}

func fmtLongDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	scales := []struct {
		Label string
		Base  time.Duration
	}{
		{"w", 7 * 24 * time.Hour},
		{"d", 24 * time.Hour},
		{"h", time.Hour},
		{"m", time.Minute},
		{"s", time.Second},
	}

	result := ""

	for _, scale := range scales {
		if d >= scale.Base {
			v := d.Truncate(scale.Base)
			r := d % scale.Base // remainder

			result = fmt.Sprintf("%s%d%s", result, int64(v/scale.Base), scale.Label)
			d = r
		}
	}

	return result
}

func parseLongDuration(s string) (time.Duration, error) {
	orig := s
	var d int64

	if s == "0" {
		return time.Duration(0), nil
	}
	if s == "" {
		return time.Duration(0), fmt.Errorf("time: invalid duration " + orig)
	}
	for s != "" {
		n := ""
		for len(s) > 0 {
			if s[0] < '0' || s[0] > '9' {
				break
			}

			n = fmt.Sprintf("%s%c", n, s[0])
			s = s[1:]
		}

		v, err := strconv.ParseInt(n, 10, 32)
		if err != nil {
			return time.Duration(0), fmt.Errorf("time: invalid duration " + orig)
		}

		switch s[0] {
		case 'w':
			d += v * 7 * 24 * int64(time.Hour)
		case 'd':
			d += v * 24 * int64(time.Hour)
		case 'h':
			d += v * int64(time.Hour)
		case 'm':
			d += v * int64(time.Minute)
		case 's':
			d += v * int64(time.Second)
		default:
			return time.Duration(0), fmt.Errorf("time: invalid duration " + orig)
		}

		s = s[1:]
	}

	return time.Duration(d), nil
}
