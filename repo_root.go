package build

import (
	"fmt"
	"path/filepath"
)

func RepoRoot() Path {
	root := FindFile(".git")
	if root == "" {
		Throw(fmt.Errorf(".git not found"))
	}
	return Path(filepath.Dir(root))
}
