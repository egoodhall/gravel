package cli

import (
	"context"
	"encoding/json"
	"os"

	"github.com/egoodhall/gravel/internal/build"
	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

type buildCmd struct {
	installFlags
	buildFlags
}

func (cmd *buildCmd) Run(ctx context.Context) error {
	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	graph, err := resolve.DependencyGraph(ctx, paths)
	if err != nil {
		return err
	}

	platforms, err := build.PlatformsFromOptions(cmd.Cross)
	if err != nil {
		return err
	}

	plan, err := build.NewPlan(paths, graph, cmd.Targets, platforms)
	if err != nil {
		return err
	}

	buildCfg := build.Config{
		Action: cmd.buildAction,
		Paths:  paths,
		Plan:   plan,
		Graph:  graph,
		Options: build.Options{
			Test:         build.TestOptions{Enabled: cmd.Test},
			Binary:       build.BinaryOptions{Enabled: cmd.Binary},
			Version:      cmd.Version,
			CrossTargets: cmd.Cross,
		},
	}

	return build.Exec(ctx, buildCfg)
}

func pkgPaths(pkgs []resolve.Pkg) []string {
	s := make([]string, len(pkgs))
	for i, p := range pkgs {
		s[i] = p.PkgPath
	}
	return s
}

func pkgPathsFromTargets(tgts []build.Target) []string {
	s := make([]string, len(tgts))
	for i, t := range tgts {
		s[i] = t.PkgPath
	}
	return s
}

func printJson(obj any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(obj)
}
