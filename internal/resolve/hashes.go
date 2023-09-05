package resolve

import (
	"github.com/egoodhall/gravel/internal/semver"
	"github.com/egoodhall/gravel/internal/types"
)

type Hashes struct {
	Old *CacheFile `json:"old"`
	New *CacheFile `json:"new"`
}

type CacheFile struct {
	Packages map[string]string `json:"packages"`
}

func (cf CacheFile) ReHash(pkg Pkg, version semver.Version) error {
	hash, err := pkg.Hash()
	if err != nil {
		return err
	}
	cf.Packages[pkg.PkgPath] = hash
	return nil
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
