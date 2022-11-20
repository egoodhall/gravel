package semver

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownSegment = errors.New("unknown segment")

//go:generate go run golang.org/x/tools/cmd/stringer -type=Segment -linecomment
type Segment byte

const (
	SegmentPatch Segment = iota // patch
	SegmentMinor                // minor
	SegmentMajor                // major
)

func (seg *Segment) UnmarshalText(p []byte) error {
	switch strings.ToLower(string(p)) {
	case "major":
		*seg = SegmentMajor
	case "minor":
		*seg = SegmentMinor
	case "patch":
		*seg = SegmentPatch
	default:
		return fmt.Errorf("%w: %s", ErrUnknownSegment, string(p))
	}
	return nil
}

func (seg Segment) MarshalText() ([]byte, error) {
	return []byte(seg.String()), nil
}
