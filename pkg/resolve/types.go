package resolve

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Pkg struct {
	Binary  string `json:"binary"`
	PkgName string `json:"pkgName"`
	PkgPath string `json:"pkgPath"`
	DirPath string `json:"dirPath"`
}

func NewPkg(frm *packages.Package) (Pkg, error) {
	dirPath := ""
	if frm.Module != nil {
		dirPath = filepath.Join(frm.Module.Dir, strings.TrimPrefix(frm.PkgPath, frm.Module.Path))
	}

	return Pkg{
		Binary:  path.Base(frm.PkgPath),
		PkgName: frm.Name,
		PkgPath: frm.PkgPath,
		DirPath: dirPath,
	}, nil
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
		switch filepath.Ext(de.Name()) {
		default:
			continue
		case ".go", ".json", ".proto", ".yml", ".yaml":
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
