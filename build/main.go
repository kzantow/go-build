package main

import (
	"fmt"
	"time"

	. "github.com/kzantow/go-build"
)

func main() {
	defer Handle()

	RunBinny()

	Cd(RepoRoot())
	Tasks(
		Task{
			Name: "format",
			Desc: "format all source files",
			Run: func() {
				//Run(`echo {{title}} {{now}}`)
				Log(`dff {{title}} {{now}}`)
				Run(`gh --version`)
				Run(`gofmt -w -s .`)
				//Run(ToolPath("gosimports"), "-local", "github.com/anchore", "-w", ".")
				//Run(`{{ToolDir}}/gosimports -local github.com/anchore -w .`)
				//Run(`{{ToolPath "gosimports"}} -local github.com/anchore -w .`)
				Run(`gosimports -local github.com/anchore -w .`)
				Run(`go mod tidy`)
			},
		},
		Task{
			Name: "lint-fix",
			Desc: "format and run lint checks",
			Deps: All("format"),
			Run: func() {
				Run("{{ToolDir}}/golangci-lint run --tests=true --fix")
				Log("lint passed!")
			},
		},
	)

	//RunTools()
}

func init() {
	Globals["title"] = "building at =>> "
	Globals["now"] = func() string {
		return fmt.Sprintf("%v", time.Now())
	}
}
