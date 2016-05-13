package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/libgit2/git2go.v24"
)

func TempGitInitPath(path string, t *testing.T) func() {
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.IsDir() == false {
		t.Fatal(err)
	}

	// Init repo
	repo, err := git.InitRepository(path, false)
	if err != nil {
		t.Fatal(err)
	}

	// Add files to index
	idx, err := repo.Index()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	err = idx.AddAll([]string{}, git.IndexAddDefault, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = idx.Write()
	if err != nil {
		t.Fatal(err)
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		t.Fatal(err)
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		t.Fatal(err)
	}

	// Initial commit
	sig := &git.Signature{
		Name:  "Test Example",
		Email: "tes@example.com",
		When:  time.Date(2016, 01, 01, 12, 00, 00, 0, time.UTC),
	}

	repo.CreateCommit("HEAD", sig, sig, "Initial commit", tree)

	// Save commmit ref for reset later
	h, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}

	obj, err := repo.Lookup(h.Target())
	if err != nil {
		t.Fatal(err)
	}

	initialCommit, err := obj.AsCommit()
	if err != nil {
		t.Fatal(err)
	}

	// Reset to initial commit, and then remove .git
	var cleanup = func() {
		repo.ResetToCommit(initialCommit, git.ResetHard, &git.CheckoutOpts{
			Strategy: git.CheckoutForce,
		})

		dotgit := filepath.Join(repo.Path())
		os.RemoveAll(dotgit)
	}

	return cleanup
}