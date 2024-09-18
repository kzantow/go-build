package main

import (
	. "github.com/kzantow/go-build"
	"github.com/kzantow/go-build/tasks"
)

func main() {
	TaskMain(
		tasks.Format,
		tasks.LintFix,
		tasks.StaticAnalysis,
		tasks.UnitTest,
		tasks.TestAll,
	)
}
