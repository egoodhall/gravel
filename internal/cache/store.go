package cache

import (
	"encoding/json"
	"os"
)

func Store(build Build, hashes Hashes) error {
	if err := os.MkdirAll(build.Paths.BinDir, 0777); err != nil {
		return err
	}

	data, err := json.MarshalIndent(build, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(build.Paths.BuildFile, data, 0666); err != nil {
		return err
	}

	data, err = json.MarshalIndent(hashes.New, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(build.Paths.HashesFile, data, 0666); err != nil {
		return err
	}

	return nil
}
