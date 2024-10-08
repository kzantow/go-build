package build

import (
	"os"
	"path/filepath"
)

// Path is used in functions that expect filepaths
type Path string

func IsDir(dir Path) bool {
	s, err := os.Stat(string(dir))
	if err != nil || s == nil {
		return false
	}
	return s.IsDir()
}

func IsRegularFile(name Path) bool {
	s, err := os.Lstat(string(name))
	if err != nil {
		return false
	}
	return !s.IsDir() && s.Mode()&os.ModeSymlink == 0
}

func FileExists(file Path) bool {
	_, err := os.Stat(string(file))
	return err == nil
}

func FindFile(glob string) string {
	dir := Get(os.Getwd())
	return findFile(dir, glob)
}

func findFile(dir string, glob string) string {
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

func PathJoin(paths ...Path) Path {
	return Path(filepath.Join(pathStrings(paths...)...))
}

func pathStrings[From string | Path](from ...From) []string {
	var out []string
	for _, v := range from {
		out = append(out, string(v))
	}
	return out
}
