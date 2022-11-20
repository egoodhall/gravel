package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"github.com/emm035/gravel/internal/semver"
)

type BumpType byte

type BumpCmd struct {
	Binary  string         `arg:"binary" required:""`
	Root    string         `name:"root" default:"."`
	Segment semver.Segment `name:"segment" default:"patch"`
	Extra   string         `name:"extra" default:""`
}

func (cmd *BumpCmd) Run() error {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, os.Kill)
	defer cancel()

	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	graph, err := resolve.DependencyGraph(ctx, paths)
	if err != nil {
		return err
	}

	pkgs := graph.Nodes().Filter(func(p resolve.Pkg) bool {
		return p.PkgName == "main" && p.Binary == cmd.Binary
	})

	for pkg := range pkgs {
		bfc, err := resolve.BuildFile(pkg)
		if err != nil {
			return err
		}

		bfc.Version.Bump(cmd.Segment, cmd.Extra)

		if err := bfc.Save(); err != nil {
			return err
		}
	}

	return nil
}
