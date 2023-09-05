package cache

import (
	"encoding/json"
	"os"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

func Store(paths gravel.Paths, hashes resolve.Hashes) error {
	if err := os.MkdirAll(paths.BinDir, 0777); err != nil {
		return err
	}

	if err := storeData(paths.HashesFile, hashes.New); err != nil {
		return err
	}

	return nil
}

func storeData(path string, data any) error {
	jsn, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, jsn, 0666); err != nil {
		return err
	}

	return nil
}
