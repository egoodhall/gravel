package resolve

import (
	"context"
	"strings"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/types"
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
		if len(pkg.Imports) == 0 {
			graph.PutNode(NewPkg(pkg))
		}

		for _, dep := range pkg.Imports {
			dpkg := NewPkg(dep)

			// Only include packages that exist within the root dir
			if !strings.HasPrefix(dpkg.DirPath, paths.RootDir) {
				continue
			}

			graph.PutEdge(NewPkg(pkg), dpkg)
		}
	}

	return graph, nil
}
