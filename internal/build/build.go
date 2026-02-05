package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/egoodhall/gravel/internal/resolve"
)

func Exec(ctx context.Context, cfg Config) error {
	if err := execTests(ctx, cfg); err != nil {
		return err
	}

	return execBinaryBuilds(ctx, cfg)
}

// parseCrossTarget splits "GOOS/GOARCH" and returns (goos, goarch, nil) or ("", "", error).
func parseCrossTarget(s string) (goos, goarch string, err error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid cross target %q: expected GOOS/GOARCH (e.g. linux/amd64)", s)
	}
	goos, goarch = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if goos == "" || goarch == "" {
		return "", "", fmt.Errorf("invalid cross target %q: GOOS and GOARCH must be non-empty", s)
	}
	if goos == "-" {
		goos = runtime.GOOS
	}
	if goarch == "-" {
		goarch = runtime.GOARCH
	}
	return goos, goarch, nil
}

// BinaryOutputPath returns the full output path for a binary built for the given platform.
// Output is always bin/GOOS_GOARCH/binary (e.g. bin/linux_amd64/myapp, bin/darwin_arm64/myapp).
func BinaryOutputPath(binDir, binary, goos, goarch string) string {
	return filepath.Join(binDir, goos+"_"+goarch, binary)
}

func generateBuildArgs(cfg Config, op BuildOp) []string {
	commit, _ := resolve.GitCommit()
	version := cfg.Options.Version

	args := []string{cfg.Action.String()}
	if cfg.Action != Install {
		// Use the absolute output path from the plan to avoid working directory issues
		args = append(args, "-o", op.OutputPath)
	}

	args = append(args,
		"-ldflags", buildLdFlags(version, commit),
		op.Target.PkgPath,
	)
	return args
}

func buildLdFlags(version, commit string) string {
	return strings.Join([]string{
		"-s -w",
		fmt.Sprintf("-X github.com/egoodhall/gravel/pkg/buildinfo.version=%s", version),
		fmt.Sprintf("-X github.com/egoodhall/gravel/pkg/buildinfo.commit=%s", commit),
	}, " ")
}
