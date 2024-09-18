package tasks

import (
	"fmt"
	"time"

	. "github.com/kzantow/go-build"
)

func init() {
	Globals["app"] = "go-build"
	Globals["now"] = func() string {
		return fmt.Sprintf("%v", time.Now())
	}
}

var Format = Task{
	Name: "format",
	Desc: "format all source files",
	Run: func() {
		Log(`Building {{app}} at {{now}}`)
		Run(`gofmt -w -s .`)
		Run(`gosimports -local github.com/anchore -w .`)
		Run(`go mod tidy`)
	},
}

var LintFix = Task{
	Name: "lint-fix",
	Desc: "format and run lint fix",
	Deps: All("format"),
	Run: func() {
		Run("golangci-lint run --fix")
		Log("lint passed!")
	},
}

var StaticAnalysis = Task{
	Name: "static-analysis",
	Desc: "run lint checks",
	Run: func() {
		Run("golangci-lint run")
		Log("lint passed!")
	},
}
