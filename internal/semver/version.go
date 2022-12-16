package semver

import (
	"fmt"
	"time"
)

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
	Extra string
}

func clone(v Version) Version {
	return Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
		Extra: v.Extra,
	}
}

// bumpSegment to the next version, setting any "lower" segments
// back to 0.
func bumpSegment(version Version, segment Segment, extra string) Version {
	version.Extra = extra
	switch segment {
	case SegmentMajor:
		version.Major++
		version.Minor = 0
		version.Patch = 0
	case SegmentMinor:
		version.Minor++
		version.Patch = 0
	case SegmentPatch:
		version.Patch++
	}
	return version
}

func bumpStrategy(version Version, strategy Strategy, extra string) Version {
	version.Extra = extra
	switch strategy {
	case StrategyDate:
		today := time.Now()
		if uint64(today.Year())%100 != version.Major {
			version.Major = uint64(today.Year()) % 100
			version.Minor = uint64(today.Month())
			version.Patch = 0
		} else if uint64(today.Month()) != version.Minor {
			version.Minor = uint64(today.Month())
			version.Patch = 0
		} else {
			version.Patch++
		}
	}
	return version
}

func (v *Version) UnmarshalText(p []byte) error {
	pv, err := Parse(string(p))
	if err != nil {
		return err
	}
	*v = *pv
	return nil
}

func (v Version) String() string {
	if v.Extra == "" {
		return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("v%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Extra)
}

func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func Parse(version string) (*Version, error) {
	v := new(Version)
	if _, err := fmt.Sscanf(version, "v%d.%d.%d-%s", &v.Major, &v.Minor, &v.Patch, &v.Extra); err == nil {
		return v, nil
	} else if _, err := fmt.Sscanf(version, "%d.%d.%d-%s", &v.Major, &v.Minor, &v.Patch, &v.Extra); err == nil {
		return v, nil
	} else if _, err := fmt.Sscanf(version, "v%d.%d.%d", &v.Major, &v.Minor, &v.Patch); err == nil {
		return v, nil
	} else if _, err := fmt.Sscanf(version, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch); err == nil {
		return v, nil
	} else {
		return nil, err
	}
}
