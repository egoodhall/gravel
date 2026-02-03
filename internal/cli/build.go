package cli

import (
	"context"
	"encoding/json"
	"os"
	"strings"

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

	crossTargets := expandCrossTargets(cmd.Cross)
	platforms, err := build.PlatformsFromOptions(crossTargets, cmd.GOOS, cmd.GOARCH)
	if err != nil {
		return err
	}
	alwaysSuffix := cmd.TargetSuffix || len(crossTargets) > 0

	plan, err := build.NewPlan(paths, graph, cmd.Targets, platforms, alwaysSuffix)
	if err != nil {
		return err
	}

	if cmd.planOnly {
		return printJson(plan)
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
			GOOS:         cmd.GOOS,
			GOARCH:       cmd.GOARCH,
			CrossTargets: expandCrossTargets(cmd.Cross),
			TargetSuffix: cmd.TargetSuffix,
		},
	}

	return build.Exec(ctx, buildCfg)
}

func expandCrossTargets(cross []string) []string {
	var out []string
	for _, s := range cross {
		for _, t := range strings.Split(s, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				out = append(out, t)
			}
		}
	}
	return out
}

func printJson(obj any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(obj)
}
