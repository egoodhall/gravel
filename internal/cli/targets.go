package cli

import (
	"context"
	"fmt"

	"github.com/egoodhall/gravel/internal/build"
	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

type targetsCmd struct {
	installFlags
	buildFlags
	Json bool `name:"json" short:"j" help:"Output the build plan as JSON (test and build package lists)."`
}

func (cmd *targetsCmd) Run(ctx context.Context) error {
	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	graph, err := resolve.DependencyGraph(ctx, paths)
	if err != nil {
		return err
	}

	if cmd.Json {
		platforms, err := build.PlatformsFromOptions(cmd.Cross)
		if err != nil {
			return err
		}
		plan, err := build.NewPlan(paths, graph, cmd.Targets, platforms)
		if err != nil {
			return err
		}
		out := struct {
			Root  string      `json:"root"`
			Test  []testEntry `json:"test"`
			Build []buildEntry `json:"build"`
		}{
			Root:  paths.RootDir,
			Test:  testEntries(plan.Test),
			Build: buildEntries(plan.Build),
		}
		return printJson(out)
	}

	for pkg := range graph.Nodes() {
		if pkg.PkgName == "main" {
			fmt.Println(pkg.PkgPath)
		}
	}

	return nil
}

type testEntry struct {
	Package string `json:"package"`
	Dir     string `json:"dir"`
}

type buildEntry struct {
	Package string `json:"package"`
	Dir     string `json:"dir"`
	Binary  string `json:"binary"`
}

func testEntries(pkgs []resolve.Pkg) []testEntry {
	entries := make([]testEntry, len(pkgs))
	for i, p := range pkgs {
		entries[i] = testEntry{
			Package: p.PkgPath,
			Dir:     p.DirPath,
		}
	}
	return entries
}

func buildEntries(ops []build.BuildOp) []buildEntry {
	entries := make([]buildEntry, len(ops))
	for i, op := range ops {
		entries[i] = buildEntry{
			Package: op.Target.PkgPath,
			Dir:     op.Target.DirPath,
			Binary:  op.OutputPath,
		}
	}
	return entries
}
