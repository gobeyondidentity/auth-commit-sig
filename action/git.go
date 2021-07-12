package action

import (
	"fmt"
	"io/ioutil"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetCommit opens the repository at repoPath and returns the commit object that
// ref resolves to.
func GetCommit(repoPath string, ref string) (*object.Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	h, err := repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve ref: %w", err)
	}

	commit, err := repo.CommitObject(*h)
	if err != nil {
		return nil, fmt.Errorf("failed to open commit: %w", err)
	}

	return commit, nil
}

// PrettyPrintCommit returns a full representation of the commit object.
func PrettyPrintCommit(commit *object.Commit) string {
	encoded := &plumbing.MemoryObject{}
	err := commit.Encode(encoded)
	if err != nil {
		panic(err)
	}

	r, err := encoded.Reader()
	if err != nil {
		panic(err)
	}

	bs, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return string(bs)
}

// EncodedCommitWithoutSignature returns an io.Reader that contains the
// canonical encoding of the commit, for signing.
func EncodedCommitWithoutSignature(commit *object.Commit) (string, error) {
	obj := &plumbing.MemoryObject{}
	err := commit.EncodeWithoutSignature(obj)
	if err != nil {
		return "", fmt.Errorf("failed to encode object: %w", err)
	}

	r, err := obj.Reader()
	if err != nil {
		return "", fmt.Errorf("failed to build reader: %w", err) // should never happen
	}

	encoded, err := ioutil.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read from encoding: %w", err) // should never happen
	}

	return string(encoded), nil
}
