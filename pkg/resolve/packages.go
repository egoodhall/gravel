package resolve

import (
	"context"
	"strings"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/pkg/types"
	"golang.org/x/tools/go/packages"
)

// DependencyGraph returns a dependency graph for the go packages
// that exist under the specified root directory. This will omit any
// package dependencies that aren't included in the root, such as
// stdlib imports or imports from other modules.
func DependencyGraph(ctx context.Context, paths gravel.Paths) (types.Graph[Pkg], error) {
	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Dir:     paths.RootDir,
		Mode:    packages.NeedDeps | packages.NeedImports | packages.NeedName | packages.NeedModule,
	}, "./...")
	if err != nil {
		return nil, err
	}

	graph := types.NewGraph[Pkg]()

	for _, pkg := range pkgs {
		pkgPkg, err := NewPkg(pkg)
		if err != nil {
			return nil, err
		}

		if len(pkg.Imports) == 0 {
			graph.PutNode(pkgPkg)
		}

		for _, dep := range pkg.Imports {
			depPkg, err := NewPkg(dep)
			if err != nil {
				return nil, err
			}

			// Only include packages that exist within the root dir
			if !strings.HasPrefix(depPkg.DirPath, paths.RootDir) {
				continue
			}

			graph.PutEdge(pkgPkg, depPkg)
		}
	}

	return graph, nil
}
