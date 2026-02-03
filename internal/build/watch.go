package build

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/fsnotify/fsnotify"
)

const watchDebounce = 150 * time.Millisecond

func Watch(ctx context.Context, paths gravel.Paths) (chan string, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := filepath.WalkDir(paths.RootDir, func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(path, paths.GravelDir) || strings.HasPrefix(path, filepath.Join(paths.RootDir, ".git")) {
			return nil
		}
		if d.IsDir() {
			return watcher.Add(filepath.Clean(path))
		}
		return nil
	}); err != nil {
		watcher.Close()
		return nil, err
	}

	events := make(chan string)

	go func() {
		defer watcher.Close()
		defer close(events)

		var (
			pending string
			timer   *time.Timer
		)
		flush := func() {
			if pending != "" {
				select {
				case events <- pending:
					pending = ""
				case <-ctx.Done():
					return
				}
			}
		}

		for {
			select {
			case event := <-watcher.Events:
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Write) || event.Has(fsnotify.Rename) {
					pending = event.Name
					if timer == nil {
						timer = time.NewTimer(watchDebounce)
					} else {
						timer.Reset(watchDebounce)
					}
				}
			case <-timer.C:
				flush()
				timer = nil
			case <-ctx.Done():
				return
			}
		}
	}()

	return events, nil
}
