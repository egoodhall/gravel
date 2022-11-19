package build

import (
	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"github.com/emm035/gravel/internal/types"
)

type Plan struct {
	Paths gravel.Paths  `json:"paths"`
	Test  []resolve.Pkg `json:"test"`
	Build []resolve.Pkg `json:"build"`
}

func NewPlan(paths gravel.Paths, graph types.Graph[resolve.Pkg], hashes resolve.Hashes) (Plan, error) {
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

	return Plan{
		Paths: paths,
		Test:  test.Slice(),
		Build: test.Filter(func(pkg resolve.Pkg) bool {
			return pkg.PkgName == "main"
		}).Slice(),
	}, nil
}
