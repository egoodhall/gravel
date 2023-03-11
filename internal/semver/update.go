package semver

import "time"

func Update(version Version, mod Mod) *Version {
	version.Extra = ""
	switch mod {
	case ModMajor:
		// Update the major version, clear the minor
		// and patch version segments
		version.Major++
		version.Minor = 0
		version.Patch = 0
	case ModMinor:
		// Update the minor version, clear the patch
		// version
		version.Minor++
		version.Patch = 0
	case ModPatch:
		// Update the patch version only
		version.Patch++
	case ModDate:
		// Update the version segments based on today's date.
		// This means that we're not actually using semantic
		// versioning. This version uses YY.MM.{build number},
		// where build number is reset whenever the year or
		// month changes.
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
	return &version
}
