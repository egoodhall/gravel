package semver

import "fmt"

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
	Extra string
}

// Bump to the next version, setting any "lower" segments
// back to 0.
func (v *Version) Bump(segment Segment, extra string) {
	v.Extra = extra
	switch segment {
	case SegmentMajor:
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case SegmentMinor:
		v.Minor++
		v.Patch = 0
	case SegmentPatch:
		v.Patch++
	}
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
