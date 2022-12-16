package build

import (
	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"github.com/emm035/gravel/internal/semver"
	"github.com/emm035/gravel/internal/types"
)

type Target struct {
	resolve.Pkg
	Version semver.Version
}

type Plan struct {
	Paths gravel.Paths  `json:"paths"`
	Test  []resolve.Pkg `json:"test"`
	Build []Target      `json:"build"`
}

func NewPlan(paths gravel.Paths, vbump semver.Bumper, graph types.Graph[resolve.Pkg], hashes resolve.Hashes) (Plan, error) {
	changedPkgPaths := hashes.ChangedPackages()
	changedPkgs := graph.Nodes().Filter(func(pkg resolve.Pkg) bool {
		return changedPkgPaths.Has(pkg.PkgPath)
	})

	dependents := graph.Transpose()
	test := types.NewSet[resolve.Pkg]()
	for pkg := range changedPkgs {
		test.Add(pkg)
		test.AddSet(dependents.Descendants(pkg))
	}

	changedVersions := hashes.ChangedVersions()
	build := types.NewSet[Target]()
	for pkg := range graph.Nodes() {
		if pkg.PkgName != "main" {
			// We only want to build main packages,
			// so this one should be skipped.
			continue
		}
		if changedVersions.Has(pkg.PkgPath) || test.Has(pkg) {
			// Either the version on this package changed (even
			// if there were no changes to the code) or an upstream
			// dependency changed, so we need to rebuild this package.
			build.Add(Target{
				Pkg:     pkg,
				Version: vbump.Bump(hashes.New.Versions[pkg.PkgPath]),
			})
			continue
		}
	}

	return Plan{
		Paths: paths,
		Test:  test.Slice(),
		Build: build.Slice(),
	}, nil
}
