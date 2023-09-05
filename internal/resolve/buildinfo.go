package resolve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/egoodhall/gravel/internal/semver"
)

var ErrBuildFileNotExist = fmt.Errorf("load build file: %w", os.ErrNotExist)

type VersionFile struct {
	semver.Version
	path string
}

func (vf *VersionFile) Save() error {
	return os.WriteFile(vf.path, []byte(vf.Version.String()), 0o660)
}

func Version(pkg Pkg) *VersionFile {
	vfpath := filepath.Join(pkg.DirPath, "version")

	vfc, err := os.ReadFile(vfpath)
	if err != nil {
		return &VersionFile{
			path: vfpath,
		}
	}

	ver := new(semver.Version)
	if err := ver.UnmarshalText(vfc); err != nil {
		return &VersionFile{
			path: vfpath,
		}
	}

	return &VersionFile{
		path:    vfpath,
		Version: *ver,
	}
}
