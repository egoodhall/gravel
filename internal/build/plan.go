package build

import (
	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"github.com/emm035/gravel/internal/types"
)

type Plan struct {
	Paths   gravel.Paths           `json:"paths"`
	Test    []resolve.Pkg          `json:"test"`
	Build   []resolve.Pkg          `json:"build"`
	Release []resolve.VersionedPkg `json:"release"`
}

func NewPlan(paths gravel.Paths, graph types.Graph[resolve.Pkg], hashes resolve.Hashes) (Plan, error) {
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
	release := types.NewSet[resolve.VersionedPkg]()
	for pkg := range graph.Nodes() {
		if changedVersions.Has(pkg.PkgPath) {
			release.Add(resolve.VersionedPkg{
				Version: hashes.New.Versions[pkg.PkgPath],
				Pkg:     pkg,
			})
		}
	}

	return Plan{
		Paths: paths,
		Test:  test.Slice(),
		Build: test.Filter(func(pkg resolve.Pkg) bool {
			return pkg.PkgName == "main"
		}).Slice(),
		Release: release.Slice(),
	}, nil
}
