package cache

import (
	"compress/gzip"
	"encoding/json"
	"os"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

func Write(paths gravel.Paths, hashes resolve.Hashes) error {
	if err := os.MkdirAll(paths.BinDir, 0777); err != nil {
		return err
	}

	if err := storeData(paths.HashesFile, hashes.New); err != nil {
		return err
	}

	return nil
}

func storeData(path string, data any) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer file.Close()

	encw := gzip.NewWriter(file)
	defer encw.Close()

	return json.NewEncoder(encw).Encode(data)
}
