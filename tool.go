package build

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"text/template"
)

var ToolDir = "{{RootDir}}/.tool"

type ToolContext struct {
	RootDir string
	ToolDir string
}

type Context map[string]any

func (c *Context) Append(key string, value any) {
	(*c)[key] = value
}

var Globals = Context{}

func init() {
	Globals.Append("RootDir", func() Path {
		return RootDir
	})
	Globals.Append("ToolDir", func() Path {
		return Path(Tpl(ToolDir))
	})
	Globals.Append("ToolPath", ToolPath)
}

func RunTools() {
	defer Handle()
	RunBinny()
	RunGoTask()
}

func RunGoTask() {
	defer appendStackOnPanic()
	if FileExists(ToolPath("task")) {
		Cd(RootDir)
		NoErr(Exec(ToolPath("task"), ExecArgs(os.Args[1:]...), ExecStd()))
	}
}

func ToolPath(toolName string) Path {
	toolPath := toolName
	switch runtime.GOOS {
	case "windows":
		toolPath += ".exe"
	}
	p := filepath.Join(Tpl(ToolDir), toolPath)
	return Path(p)
}

func Tpl(template string, args ...map[string]any) string {
	context := map[string]any{}
	for k, v := range Globals {
		if reflect.TypeOf(v).Kind() != reflect.Func {
			context[k] = v
		}
	}
	for _, arg := range args {
		for k, v := range arg {
			context[k] = v
		}
	}
	return render(template, context)
}

func render(tpl string, context map[string]any) string {
	funcs := template.FuncMap{}
	for k, v := range Globals {
		v := v
		val := reflect.ValueOf(v)
		switch val.Type().Kind() {
		case reflect.Func:
			funcs[k] = v
		case reflect.String:
			funcs[k] = func() string { return Tpl(val.String()) }
		default:
			funcs[k] = func() any { return v }
		}
	}
	t := Get(template.New("test").Funcs(funcs).Parse(tpl))
	var buf bytes.Buffer
	NoErr(t.Execute(&buf, context))
	return buf.String()
}
