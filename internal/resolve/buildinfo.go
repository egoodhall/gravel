package resolve

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ErrBuildFileNotExist = fmt.Errorf("load build file: %w", os.ErrNotExist)

var buildFileNames = []string{
	".gravel.yml",
	".gravel.yaml",
}

type BuildConfig struct {
	Version string `yaml:"version"`
}

func BuildFile(pkg Pkg) (*BuildConfig, error) {
	for _, fileName := range buildFileNames {
		data, err := os.ReadFile(filepath.Join(pkg.DirPath, fileName))
		if os.IsNotExist(err) {
			fmt.Println(filepath.Join(pkg.DirPath, fileName))
			continue
		} else if err != nil {
			return nil, err
		}

		bf := new(BuildConfig)
		if err := yaml.Unmarshal(data, bf); err != nil {
			return nil, err
		}

		return bf, nil
	}
	return nil, ErrBuildFileNotExist
}
