package tools

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// SampleSpec describes how backups should be sampled
// for a given time range
type SampleSpec struct {
	Before   time.Duration
	Interval time.Duration
	Round    time.Duration
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
	if s.Before == 0 {
		return true
	}

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

	if !t.Truncate(s.Interval).Equal(t.Round(s.Round)) {
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
	out := ""

	if s.Before > 0 {
		out = fmt.Sprintf("%s@%s", out, fmtLongDuration(s.Before))
	}

	if s.Interval > 0 {
		out = fmt.Sprintf("%s/%s", out, fmtLongDuration(s.Interval))
	}

	if s.Round > 0 {
		out = fmt.Sprintf("%s~%s", out, fmtLongDuration(s.Round))
	}

	return out, nil
}

func (s *SampleSpec) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	return s.UnmarshalString(str)
}

func (s *SampleSpec) UnmarshalString(str string) error {
	orig := str

	for len(str) > 0 {
		switch str[0] {
		case '@': // Before
			head, tail := extractLongDuration(str[1:])
			str = tail

			d, err := parseLongDuration(head)
			if err != nil {
				return err
			}

			s.Before = d
		case '/': // Interval
			head, tail := extractLongDuration(str[1:])
			str = tail

			d, err := parseLongDuration(head)
			if err != nil {
				return err
			}

			s.Interval = d
		case '~': // Round
			head, tail := extractLongDuration(str[1:])
			str = tail

			d, err := parseLongDuration(head)
			if err != nil {
				return err
			}

			s.Round = d
		default:
			return fmt.Errorf("tools: invalid sample spec %s: unrecognized control character '%c'", orig, str[0])
		}
	}

	return nil
}

var longDurationTimeScales = []struct {
	Label byte
	Base  time.Duration
}{
	{'w', 7 * 24 * time.Hour},
	{'d', 24 * time.Hour},
	{'h', time.Hour},
	{'m', time.Minute},
	{'s', time.Second},
}

func fmtLongDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	result := ""

	for _, scale := range longDurationTimeScales {
		if d >= scale.Base {
			v := d.Truncate(scale.Base)
			r := d % scale.Base // remainder

			result = fmt.Sprintf("%s%d%c", result, int64(v/scale.Base), scale.Label)
			d = r
		}
	}

	return result
}

func extractLongDuration(s string) (head, tail string) {
	for len(s) > 0 {
		if s[0] >= '0' && s[0] <= '9' {
			head = fmt.Sprintf("%s%c", head, s[0])
			s = s[1:]
			continue
		}

		matched := false
		for _, scale := range longDurationTimeScales {
			if s[0] == scale.Label {
				head = fmt.Sprintf("%s%c", head, s[0])
				s = s[1:]
				matched = true
				break
			}
		}

		if matched {
			continue
		} else {
			tail = s
			return
		}
	}

	tail = ""
	return
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

		matched := false
		for _, scale := range longDurationTimeScales {
			if s[0] == scale.Label {
				d += v * int64(scale.Base)
				matched = true
				break
			}
		}

		if !matched {
			return time.Duration(0), fmt.Errorf("time: invalid duration " + orig)
		}

		s = s[1:]
	}

	return time.Duration(d), nil
}
