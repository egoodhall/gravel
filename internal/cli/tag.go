package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
	"github.com/egoodhall/gravel/internal/semver"
)

type releaseCmd struct {
	installFlags
	Mods map[string]semver.Mod `name:"mod" help:"Modifications to apply to existing versions"`
	Tag  bool                  `name:"tag" help:"Write the updated versions as git tags"`
}

func (cmd *releaseCmd) Run() error {
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

	mainPkgs := graph.Nodes().Filter(func(p resolve.Pkg) bool {
		return p.PkgName == "main"
	})

	versions, err := semver.LoadTags(paths)
	if err != nil {
		return err
	}

	nextVersions := make(map[string]*semver.Version)
	for pkg := range mainPkgs {
		if mod, ok := cmd.Mods[pkg.Binary]; ok {
			nextVersions[pkg.Binary] = semver.Update(*versions[pkg.Binary], mod)
		}
	}

	if !cmd.Tag {
		maxlen := 0
		for bin := range nextVersions {
			if len(bin) > maxlen {
				maxlen = len(bin)
			}
		}

		for binary, nextVersion := range nextVersions {
			if version, ok := versions[binary]; ok {
				version.Extra = ""
				fmt.Printf(fmt.Sprintf("%%%ds : %%s -> %%s\n", maxlen), binary, version, nextVersion)
			}
		}

		return nil
	}

	return semver.WriteTags(paths, nextVersions)
}
