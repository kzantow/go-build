package build

import "os"

var RootDir = RepoRoot()

func Cd[path string | Path](dir path) {
	NoErr(os.Chdir(string(dir)))
}

func Cwd() Path {
	return Path(Get(os.Getwd()))
}
