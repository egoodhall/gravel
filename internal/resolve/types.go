package resolve

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emm035/gravel/internal/semver"
	"golang.org/x/tools/go/packages"
)

type Module struct {
	ModPath string `json:"modPath"`
	DirPath string `json:"dirPath"`
}

type VersionedPkg struct {
	Pkg
	Version semver.Version `json:"version"`
}

type Pkg struct {
	Module  Module `json:"module"`
	Binary  string `json:"binary"`
	PkgName string `json:"pkgName"`
	PkgPath string `json:"pkgPath"`
	DirPath string `json:"dirPath"`
}

func NewPkg(frm *packages.Package) (pkg Pkg) {
	pkg.PkgName = frm.Name
	pkg.PkgPath = frm.PkgPath
	pkg.Binary = path.Base(pkg.PkgPath)
	if frm.Module != nil {
		pkg.DirPath = filepath.Join(frm.Module.Dir, strings.TrimPrefix(frm.PkgPath, frm.Module.Path))
		pkg.Module.DirPath = frm.Module.Dir
		pkg.Module.ModPath = frm.Module.Path
	}
	return
}

func (pkg Pkg) String() string {
	return pkg.PkgPath
}

func (pkg Pkg) Hash() (string, error) {
	hash := sha256.New()

	des, err := os.ReadDir(pkg.DirPath)
	if err != nil {
		return "", err
	}

	for _, de := range des {
		if de.IsDir() {
			continue
		}

		f, err := os.Open(filepath.Join(pkg.DirPath, de.Name()))
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(hash, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
