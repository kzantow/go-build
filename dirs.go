package build

import "os"

var RootDir = func() Path {
	return Path(Tpl("{{RepoRoot}}"))
}

func Cd[path string | Path](dir path) {
	NoErr(os.Chdir(string(dir)))
}

func Cwd() Path {
	return Path(Get(os.Getwd()))
}
