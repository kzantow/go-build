package tasks

import . "github.com/kzantow/go-build"

var UnitTest = Task{
	Name: "unit",
	Desc: "run unit tests",
	Run: func() {
		Run(`go test ./...`)
	},
}

var TestAll = Task{
	Name: "test",
	Deps: All("unit"),
}
