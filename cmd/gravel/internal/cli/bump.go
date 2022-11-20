package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"github.com/emm035/gravel/internal/semver"
	"github.com/go-git/go-git/v5"
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

		if err := cmd.Tag(pkg, bfc.Version); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *BumpCmd) Tag(pkg resolve.Pkg, version semver.Version) error {
	repo, err := git.PlainOpen(cmd.Root)
	if err != nil {
		return err
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	fmt.Printf("%s/%s", pkg.Binary, version)
	_, err = repo.CreateTag(fmt.Sprintf("%s/%s", pkg.Binary, version), ref.Hash(), nil)
	return err
}
