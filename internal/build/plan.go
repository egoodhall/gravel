package build

import (
	"os"
	"path/filepath"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
	"github.com/egoodhall/gravel/internal/semver"
	"github.com/egoodhall/gravel/internal/types"
)

type Target struct {
	resolve.Pkg
	Version *semver.Version
}

type Plan struct {
	Paths gravel.Paths  `json:"paths"`
	Test  []resolve.Pkg `json:"test"`
	Build []Target      `json:"build"`
}

func NewPlan(paths gravel.Paths, graph types.Graph[resolve.Pkg], hashes resolve.Hashes, targets []string) (Plan, error) {
	changedPkgPaths := hashes.ChangedPackages()
	changedPkgs := graph.Nodes().Filter(func(pkg resolve.Pkg) bool {
		return changedPkgPaths.Has(pkg.PkgPath)
	})

	dependents := graph.Transpose()
	testPkgs := types.NewSet[resolve.Pkg]()
	for pkg := range changedPkgs {
		testPkgs.Add(pkg)
		testPkgs.AddSet(dependents.Descendants(pkg))
	}

	targetPaths := findDirs(targets)
	buildPkgs := types.NewSet[resolve.Pkg]()
	for pkg := range graph.Nodes() {
		if pkg.PkgName != "main" {
			// We only want to build main packages,
			// so this one should be skipped.
			continue
		}

		if (len(targets) == 0 && testPkgs.Has(pkg)) || targetPaths.Has(pkg.DirPath) || matchesSpecificTarget(targets, pkg) {
			// Either an upstream dependency changed,
			// so we need to rebuild this package,
			// or it matched a specified target
			buildPkgs.Add(pkg)
			continue
		}
	}

	versions, err := semver.LoadTags(paths)
	if err != nil {
		return Plan{}, err
	}

	buildTgts := types.NewSet[Target]()
	for pkg := range buildPkgs {
		tgt := Target{
			Pkg:     pkg,
			Version: semver.Zero(),
		}
		if version, ok := versions[tgt.Binary]; ok {
			tgt.Version = version
		}
		buildTgts.Add(tgt)
	}

	return Plan{
		Paths: paths,
		Test:  testPkgs.Slice(),
		Build: buildTgts.Slice(),
	}, nil
}

func findDirs(targets []string) types.Set[string] {
	found := types.NewSet[string]()
	for _, target := range targets {
		matches, err := filepath.Glob(target)
		if err != nil {
			continue
		}

		for _, match := range matches {
			abspath, err := filepath.Abs(match)
			if err != nil {
				continue
			}
			fi, err := os.Stat(abspath)
			if err != nil {
				continue
			}
			if fi.IsDir() {
				found.Add(abspath)
			}
		}
	}
	return found
}

func matchesSpecificTarget(targets []string, pkg resolve.Pkg) bool {
	for _, target := range targets {
		if target == pkg.PkgPath || target == pkg.Binary {
			return true
		}
	}
	return false
}
