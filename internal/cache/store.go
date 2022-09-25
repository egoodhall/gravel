package cache

import (
	"encoding/json"
	"os"
)

func Store(build Build, hashes Hashes, planOnly bool) error {
	if err := os.MkdirAll(build.Paths.BinDir, 0777); err != nil {
		return err
	}

	if err := storeData(build.Paths.BuildFile, build); err != nil {
		return err
	}

	if !planOnly {
		if err := storeData(build.Paths.HashesFile, hashes.New); err != nil {
			return err
		}
	}

	return nil
}

func storeData(path string, data any) error {
	jdata, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, jdata, 0666); err != nil {
		return err
	}

	return nil
}
