package tools

import (
	"sort"
	"time"
)

// Matches determines whether a given time matches any number
// of SampleSpecs provided by ensuring that:
//
// - The time is outside the SampleSpec's window
//
// --OR--
//
// - The time is within a SampleSpec's window and matches its interval
func Matches(t time.Time, specs []SampleSpec) bool {
	sort.Sort(byAge(specs))

	for _, spec := range specs {
		if spec.WithinWindow(t) {
			return spec.Matches(t.Round(spec.Round))
		}
	}

	return true
}

type byAge []SampleSpec

func (s byAge) Len() int { return len(s) }

func (s byAge) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s byAge) Less(i, j int) bool { return s[i].Before > s[j].Before }
