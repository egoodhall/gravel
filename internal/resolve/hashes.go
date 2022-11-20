package resolve

import (
	"github.com/emm035/gravel/internal/semver"
	"github.com/emm035/gravel/internal/types"
)

type Hashes struct {
	Old *CacheFile
	New *CacheFile
}

type CacheFile struct {
	Packages map[string]string         `json:"packages"`
	Versions map[string]semver.Version `json:"versions"`
}

func (h Hashes) ChangedPackages() types.Set[string] {
	s := types.NewSet[string]()
	for pkg, hash := range h.New.Packages {
		if hash != h.Old.Packages[pkg] {
			s.Add(pkg)
		}
	}
	for pkg, hash := range h.Old.Packages {
		if hash != h.New.Packages[pkg] {
			s.Add(pkg)
		}
	}
	return s
}

func (h Hashes) ChangedVersions() types.Set[string] {
	s := types.NewSet[string]()
	for pkg, hash := range h.New.Versions {
		if hash != h.Old.Versions[pkg] {
			s.Add(pkg)
		}
	}
	for pkg, hash := range h.Old.Versions {
		if hash != h.New.Versions[pkg] {
			s.Add(pkg)
		}
	}
	return s
}
