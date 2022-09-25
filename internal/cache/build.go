package cache

import (
	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/pkg/resolve"
	"github.com/emm035/gravel/pkg/types"
)

type Build struct {
	Paths gravel.Paths  `json:"paths"`
	Test  []resolve.Pkg `json:"test"`
	Build []resolve.Pkg `json:"build"`
}

func NewBuild(paths gravel.Paths, graph types.Graph[resolve.Pkg], hashes Hashes) (Build, error) {
	changedPkgPaths := hashes.Changed()
	changed := graph.Nodes().Filter(func(pkg resolve.Pkg) bool {
		return changedPkgPaths.Has(pkg.PkgPath)
	})

	dependents := graph.Transpose()
	test := types.NewSet[resolve.Pkg]()
	for pkg := range changed {
		test.Add(pkg)
		test.AddSet(dependents.Descendants(pkg))
	}

	return Build{
		Paths: paths,
		Test:  test.Slice(),
		Build: test.Filter(func(pkg resolve.Pkg) bool {
			return pkg.PkgName == "main"
		}).Slice(),
	}, nil
}
