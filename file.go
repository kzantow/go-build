package build

import (
	"os"
	"path/filepath"
)

func IsDir(dir string) bool {
	s, err := os.Stat(dir)
	if err != nil || s == nil {
		return false
	}
	return s.IsDir()
}

func IsRegularFile(name string) bool {
	s, err := os.Lstat(name)
	if err != nil {
		return false
	}
	return !s.IsDir() && s.Mode()&os.ModeSymlink == 0
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func PathJoin(path ...string) string {
	return filepath.Join(path...)
}
