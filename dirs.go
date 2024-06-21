package build

import "os"

var RootDir = RepoRoot()

func Cd(dir string) {
	NoErr(os.Chdir(dir))
}

func Cwd() string {
	return Get(os.Getwd())
}
