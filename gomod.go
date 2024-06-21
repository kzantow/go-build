package build

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

func GoDepVersion(module string) string {
	modFile := FindFile("go.mod")
	contents := Get(os.ReadFile(modFile))
	f := Get(modfile.Parse("go.mod", contents, nil))
	if f.Module != nil && f.Module.Mod.Path == module {
		return GitRevision()
	}
	for _, r := range f.Require {
		if r.Mod.Path == module {
			return r.Mod.Version
		}
	}
	return "UNKNOWN"
}

func GitRevision() string {
	buf := bytes.Buffer{}
	err := Exec("git", ExecArgs("rev-parse", "--short", "HEAD"), ExecOut(&buf))
	if err != nil {
		return "UNKNOWN"
	}
	return strings.TrimSpace(buf.String())
}

func FindFile(glob string) string {
	dir := Get(os.Getwd())
	for {
		f := filepath.Join(dir, glob)
		matches, _ := filepath.Glob(f)
		if len(matches) > 0 {
			return matches[0]
		}
		if dir == filepath.Dir(dir) {
			return ""
		}
		dir = filepath.Dir(dir)
	}
}
