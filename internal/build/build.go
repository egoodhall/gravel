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

func Exec(ctx context.Context, build Plan) error {

	if len(build.Test) > 0 {
		if err := execTests(ctx, build.Test); err != nil {
			return err
		}
	}

	if len(build.Build) > 0 {
		if err := execBuilds(ctx, build.Paths, build.Build); err != nil {
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

func execBuilds(parent context.Context, paths gravel.Paths, pkgs []resolve.Pkg) error {
	eg, ctx := errgroup.WithContext(parent)
	for _, pkg := range pkgs {
		fmt.Println(pkg.PkgPath)
		eg.Go(execBuild(ctx, paths, pkg))
	}
	return eg.Wait()
}

func execBuild(ctx context.Context, paths gravel.Paths, pkg resolve.Pkg) func() error {
	var version string
	if bc, err := resolve.BuildFile(pkg); err == nil {
		version = bc.Version
	}

	commit, _ := resolve.GitCommit()

	cmd := exec.CommandContext(ctx, "go", "build",
		"-ldflags", buildLdFlags(version, commit),
		"-o", filepath.Join(paths.BinDir, pkg.Binary),
		pkg.PkgPath,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run
}

func buildLdFlags(version, commit string) string {
	return strings.Join([]string{
		"-s -w",
		fmt.Sprintf("-X github.com/emm035/gravel/pkg/buildinfo.version=%s", version),
		fmt.Sprintf("-X github.com/emm035/gravel/pkg/buildinfo.commit=%s", commit),
	}, " ")
}
