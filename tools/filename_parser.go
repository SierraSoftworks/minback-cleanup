package tools

import (
	"path"
	"time"

	"github.com/araddon/dateparse"
)

// ParseFilename extracts a timestamp from a given filename,
// attempting to intelligently detect the timestamp format
// used.
func ParseFilename(name, prefix string) (time.Time, error) {
	ext := path.Ext(name)

	datePart := name[len(prefix) : len(name)-len(ext)]
	return dateparse.ParseAny(datePart)
}
