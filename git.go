package build

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

func GitRoot() Path {
	root := FindFile(".git")
	if root == "" {
		Throw(fmt.Errorf(".git not found"))
	}
	return Path(filepath.Dir(root))
}

func GitRevision() string {
	buf := bytes.Buffer{}
	err := Exec("git", ExecArgs("rev-parse", "--short", "HEAD"), ExecOut(&buf))
	if err != nil {
		return "UNKNOWN"
	}
	return strings.TrimSpace(buf.String())
}
