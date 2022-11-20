package resolve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/emm035/gravel/internal/semver"
	"gopkg.in/yaml.v3"
)

var ErrBuildFileNotExist = fmt.Errorf("load build file: %w", os.ErrNotExist)

var buildFileNames = []string{
	".gravel.yml",
	".gravel.yaml",
}

type BuildFileContents struct {
	Path    string         `yaml:"-"`
	Version semver.Version `yaml:"version"`
}

func (bfc *BuildFileContents) Save() error {
	data, err := yaml.Marshal(bfc)
	if err != nil {
		return err
	}
	return os.WriteFile(bfc.Path, data, 0o660)
}

func BuildFile(pkg Pkg) (*BuildFileContents, error) {
	for _, fileName := range buildFileNames {
		data, err := os.ReadFile(filepath.Join(pkg.DirPath, fileName))
		if os.IsNotExist(err) {
			fmt.Println(filepath.Join(pkg.DirPath, fileName))
			continue
		} else if err != nil {
			return nil, err
		}

		bf := new(BuildFileContents)
		if err := yaml.Unmarshal(data, bf); err != nil {
			return nil, err
		}

		bf.Path = filepath.Join(pkg.DirPath, fileName)

		return bf, nil
	}
	return nil, ErrBuildFileNotExist
}
