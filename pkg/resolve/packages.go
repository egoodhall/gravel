package resolve

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/emm035/gravel/pkg/types"
	"golang.org/x/tools/go/packages"
)

func DependencyGraph(ctx context.Context, root string, queries ...string) (types.Graph[Pkg], error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Dir:     absRoot,
		Mode:    packages.NeedDeps | packages.NeedImports | packages.NeedName | packages.NeedModule,
	}, queries...)
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

			if !strings.HasPrefix(depPkg.DirPath, absRoot) {
				continue
			}

			graph.PutEdge(pkgPkg, depPkg)
		}
	}

	return graph, nil
}
