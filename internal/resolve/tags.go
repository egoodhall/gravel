package resolve

import (
	"strings"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/semver"
	"github.com/go-git/go-git/v5"
)

func AllTags(paths gravel.Paths) (map[string]*semver.Version, error) {
	repo, err := git.PlainOpen(paths.RootDir)
	if err != nil {
		return nil, err
	}

	refs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	extra, err := getExtra(repo)
	if err != nil {
		return nil, err
	}

	clean, err := isClean(repo)
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	versions := make(map[string]*semver.Version)
	for ref, err := refs.Next(); err == nil; ref, err = refs.Next() {
		if !ref.Name().IsTag() {
			continue
		}

		tag := ref.Name().Short()
		bin, ver, ok := strings.Cut(tag, "/")
		if !ok {
			continue
		}

		sver, err := semver.Parse(ver)
		if err != nil {
			return nil, err
		}

		tagCommit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return nil, err
		}

		if headCommit.Hash != tagCommit.Hash || !clean {
			sver.Extra = extra
		}

		if cver, ok := versions[bin]; ok {
			versions[bin] = semver.Max(sver, cver)
		} else {
			versions[bin] = sver
		}
	}

	return versions, nil
}

func getExtra(repo *git.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", err
	}

	var extra string
	if head.Name().IsBranch() {
		extra = head.Name().Short()
	} else {
		extra = commit.Hash.String()
	}

	return extra, nil
}

func isClean(repo *git.Repository) (bool, error) {
	tree, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := tree.Status()
	if err != nil {
		return false, err
	}

	return status.IsClean(), nil
}
