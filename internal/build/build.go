package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"golang.org/x/sync/errgroup"
)

func Exec(ctx context.Context, action Action, build Plan) error {

	if len(build.Test) > 0 {
		if err := execTests(ctx, build.Test); err != nil {
			return err
		}
	}

	if len(build.Build) > 0 {
		if err := execBuilds(ctx, action, build.Paths, build.Build); err != nil {
			return err
		}
	}

	return nil
}

func execTests(ctx context.Context, pkgs []resolve.Pkg) error {
	args := []string{"test", "-count", "1"}
	for _, pkg := range pkgs {
		args = append(args, pkg.PkgPath)
	}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func execBuilds(parent context.Context, action Action, paths gravel.Paths, tgts []Target) error {
	eg, ctx := errgroup.WithContext(parent)
	for _, tgt := range tgts {
		fmt.Println(tgt.PkgPath)
		eg.Go(execBuild(ctx, action, paths, tgt))
	}
	return eg.Wait()
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=Action -linecomment
type Action byte

const (
	Build   Action = iota // build
	Install               // install
)

func execBuild(ctx context.Context, action Action, paths gravel.Paths, tgt Target) func() error {
	commit, _ := resolve.GitCommit()

	args := []string{action.String()}
	if action != Install {
		args = append(args, "-o", filepath.Join(paths.BinDir, tgt.Binary))
	}
	args = append(args,
		"-ldflags", buildLdFlags(tgt.Version.String(), commit),
		tgt.PkgPath,
	)

	c := exec.CommandContext(ctx, "go", args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	return c.Run
}

func buildLdFlags(version, commit string) string {
	return strings.Join([]string{
		"-s -w",
		fmt.Sprintf("-X github.com/emm035/gravel/pkg/buildinfo.version=%s", version),
		fmt.Sprintf("-X github.com/emm035/gravel/pkg/buildinfo.commit=%s", commit),
	}, " ")
}
