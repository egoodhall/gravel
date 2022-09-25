package cache

import (
	"encoding/json"
	"os"

	"github.com/emm035/gravel/internal/build"
	"github.com/emm035/gravel/pkg/resolve"
)

func Store(plan build.Plan, hashes resolve.Hashes, planOnly bool) error {
	if err := os.MkdirAll(plan.Paths.BinDir, 0777); err != nil {
		return err
	}

	if err := os.WriteFile(plan.Paths.GitignoreFile, []byte{'*'}, 0666); err != nil {
		return err
	}

	if err := storeData(plan.Paths.PlanFile, plan); err != nil {
		return err
	}

	if !planOnly {
		if err := storeData(plan.Paths.HashesFile, hashes.New); err != nil {
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
