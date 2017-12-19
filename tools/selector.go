package tools

import "time"

// Matches determines whether a given time matches any number
// of SampleSpecs provided by ensuring that:
//
// - The time is outside the SampleSpec's window OR
// - The time is within a SampleSpec's window and matches its interval
func Matches(t time.Time, specs []SampleSpec) bool {
	for _, spec := range specs {
		if spec.WithinWindow(t) && !spec.Matches(t) {
			return false
		}
	}

	return true
}
