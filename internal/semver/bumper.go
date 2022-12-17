package semver

import (
	"errors"
)

type Bumper interface {
	Bump(Version) Version
}

type bumperFunc func(Version) Version

func (bf bumperFunc) Bump(v Version) Version {
	return bf(v)
}

func NewBumper(segment *Segment, strategy *Strategy, extra string) (Bumper, error) {
	if (segment != nil) && (strategy != nil) {
		return nil, errors.New("segment and strategy are not exclusive")
	}
	if segment != nil {
		return bumperFunc(func(v Version) Version {
			return bumpSegment(v, *segment, extra)
		}), nil
	} else if strategy != nil {
		// Use a standard bump strategy for incrementing
		return bumperFunc(func(v Version) Version {
			return bumpStrategy(v, *strategy, extra)
		}), nil
	} else {
		// If neither segment nor strategy was provided,
		// return a bumper that will simply copy the version
		// and set the extra value
		return bumperFunc(func(v Version) Version {
			v.Extra = extra
			return v
		}), nil
	}
}
