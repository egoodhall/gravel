package build

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
)

func Exec(ctx context.Context, cfg Config) error {
	if err := execTests(ctx, cfg); err != nil {
		return err
	}

	if err := execBinaryBuilds(ctx, cfg); err != nil {
		return err
	}

	if err := execDockerBuilds(ctx, cfg); err != nil {
		return err
	}

	return nil
}

func generateBuildArgs(action Action, paths gravel.Paths, tgt Target) []string {
	commit, _ := resolve.GitCommit()

	args := []string{action.String()}
	if action != Install {
		outPath := filepath.Join(paths.BinDir, tgt.Binary)
		if p, err := filepath.Rel(tgt.Module.DirPath, outPath); err == nil {
			outPath = "./" + p
		}
		args = append(args, "-o", outPath)
	}

	args = append(args,
		"-ldflags", buildLdFlags(tgt.Version.String(), commit),
		tgt.PkgPath,
	)
	return args
}

func buildLdFlags(version, commit string) string {
	return strings.Join([]string{
		"-s -w",
		fmt.Sprintf("-X github.com/emm035/gravel/pkg/buildinfo.version=%s", version),
		fmt.Sprintf("-X github.com/emm035/gravel/pkg/buildinfo.commit=%s", commit),
	}, " ")
}
