package resolve

import "github.com/go-git/go-git/v5"

func GitCommit() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}
