package build

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func RepoRoot() string {
	root := FindFile(".git")
	if root == "" {
		Throw(fmt.Errorf(".git not found"))
	}
	return filepath.Dir(root)
}

func RepoRootRevParse() string {
	root, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		panic(fmt.Errorf("unable to find repo root dir: %w", err))
	}
	absRepoRoot, err := filepath.Abs(strings.TrimSpace(string(root)))
	if err != nil {
		panic(fmt.Errorf("unable to get abs path to repo root: %w", err))
	}
	return absRepoRoot
}
