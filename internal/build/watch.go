package build

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
	"github.com/fsnotify/fsnotify"
)

const watchDebounce = 150 * time.Millisecond

// Watch watches only directories for packages in the given set (typically packages in the dependency graph for build targets).
func Watch(ctx context.Context, paths gravel.Paths, watchPkgs []resolve.Pkg) (chan string, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Track which directories we've added to avoid duplicates
	watchedDirs := make(map[string]bool)
	skipDir := paths.SkipDir
	gitDir := filepath.Join(paths.RootDir, ".git")

	// Watch directories for all relevant packages (excluding bin and .git)
	for _, pkg := range watchPkgs {
		dir := filepath.Clean(pkg.DirPath)
		if !watchedDirs[dir] && strings.HasPrefix(dir, paths.RootDir) {
			// Skip bin directory and .git
			if strings.HasPrefix(dir, skipDir) || strings.HasPrefix(dir, gitDir) {
				continue
			}
			if err := watcher.Add(dir); err == nil {
				watchedDirs[dir] = true
			}
		}
	}

	events := make(chan string)

	go func() {
		defer watcher.Close()
		defer close(events)

		skipDir := paths.SkipDir
		gitDir := filepath.Join(paths.RootDir, ".git")

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
			// Only select on timer.C when timer is set; nil channel never triggers.
			var timerCh <-chan time.Time
			if timer != nil {
				timerCh = timer.C
			}

			select {
			case event := <-watcher.Events:
				// Filter out events in bin directory or .git
				eventPath := filepath.Clean(event.Name)
				if strings.HasPrefix(eventPath, skipDir) || strings.HasPrefix(eventPath, gitDir) {
					continue
				}
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Write) || event.Has(fsnotify.Rename) {
					pending = event.Name
					if timer == nil {
						timer = time.NewTimer(watchDebounce)
					} else {
						timer.Reset(watchDebounce)
					}
				}
			case <-timerCh:
				flush()
				timer = nil
			case <-ctx.Done():
				return
			}
		}
	}()

	return events, nil
}
