package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/egoodhall/gravel/internal/build"
	"github.com/egoodhall/gravel/internal/gravel"
)

type watchCmd struct {
	installFlags
	buildFlags
}

func (cmd *watchCmd) Run() error {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, os.Kill)
	defer cancel()

	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	events, err := build.Watch(ctx, paths)
	if err != nil {
		return err
	}

	fmt.Printf("Watching files in %s for changes\n", paths.RootDir)

	for {
		select {
		case file := <-events:
			if relfile, err := filepath.Rel(paths.RootDir, file); err == nil {
				file = relfile
			}

			fmt.Printf("%s changed, rebuilding.\n", file)

			if err := (&buildCmd{
				installFlags: cmd.installFlags,
				buildFlags:   cmd.buildFlags,
			}).Run(); err != nil {
				fmt.Printf("Build failed: %s\n", file)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
