package build

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
	"github.com/egoodhall/gravel/internal/types"
)

// Platform is a GOOS/GOARCH build target.
type Platform struct {
	GOOS   string
	GOARCH string
}

type Target struct {
	resolve.Pkg
}

type Plan struct {
	Paths gravel.Paths  `json:"paths"`
	Test  []resolve.Pkg `json:"test"`
	Build []Target      `json:"build"`
}

// PlatformsFromOptions returns the list of platforms we are building for:
// if crossTargets is non-empty, one platform per "GOOS/GOARCH" entry; otherwise a single platform from goos/goarch (defaulting to current OS/arch).
func PlatformsFromOptions(crossTargets []string, goos, goarch string) ([]Platform, error) {
	if len(crossTargets) > 0 {
		out := make([]Platform, 0, len(crossTargets))
		for _, t := range crossTargets {
			g, a, err := parseCrossTarget(t)
			if err != nil {
				return nil, err
			}
			out = append(out, Platform{GOOS: g, GOARCH: a})
		}
		return out, nil
	}
	g, a := goos, goarch
	if g == "" {
		g = runtime.GOOS
	}
	if a == "" {
		a = runtime.GOARCH
	}
	return []Platform{{GOOS: g, GOARCH: a}}, nil
}

func NewPlan(paths gravel.Paths, graph types.Graph[resolve.Pkg], targets []string, platforms []Platform, alwaysSuffix bool) (Plan, error) {
	targetPaths := findDirs(targets)

	// Build set: all main packages if no targets specified, otherwise only those matching targets.
	buildPkgs := types.NewSet[resolve.Pkg]()
	for pkg := range graph.Nodes() {
		if pkg.PkgName != "main" {
			continue
		}
		if len(targets) == 0 {
			buildPkgs.Add(pkg)
			continue
		}
		if targetPaths.Has(pkg.DirPath) || matchesSpecificTarget(targets, pkg) {
			buildPkgs.Add(pkg)
		}
	}

	// Include any main package whose output binary is missing for any platform we're building for.
	for pkg := range graph.Nodes() {
		if pkg.PkgName != "main" {
			continue
		}
		for _, plat := range platforms {
			name := BinaryOutputName(pkg.Binary, plat.GOOS, plat.GOARCH, alwaysSuffix)
			if _, err := os.Stat(filepath.Join(paths.BinDir, name)); err != nil {
				if os.IsNotExist(err) {
					buildPkgs.Add(pkg)
					break
				}
			}
		}
	}

	// Test packages: each build target and all of its dependencies.
	testPkgs := types.NewSet[resolve.Pkg]()
	for pkg := range buildPkgs {
		testPkgs.Add(pkg)
		testPkgs.AddSet(graph.Descendants(pkg))
	}

	buildTgts := types.NewSet[Target]()
	for pkg := range buildPkgs {
		buildTgts.Add(Target{Pkg: pkg})
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
