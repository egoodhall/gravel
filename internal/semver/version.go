package semver

import (
	"fmt"
)

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
	Extra string
}

func Zero() *Version {
	return &Version{}
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

func Max(v1, v2 *Version) *Version {
	if v1 == nil && v2 == nil {
		return nil
	} else if v2 == nil {
		return v1
	} else if v1 == nil {
		return v2
	} else if v1.Major > v2.Major {
		return v1
	} else if v1.Major < v2.Major {
		return v2
	} else if v1.Minor > v2.Minor {
		return v1
	} else if v1.Minor < v2.Minor {
		return v2
	} else if v1.Patch > v2.Patch {
		return v1
	} else if v1.Patch < v2.Patch {
		return v2
	} else if v1.Extra == "" && v2.Extra != "" {
		return v1
	} else if v1.Extra != "" && v2.Extra == "" {
		return v2
	} else {
		return v2
	}
}
