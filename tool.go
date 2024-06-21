package build

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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
	Globals.Append("RootDir", func() string {
		return RootDir
	})
	Globals.Append("ToolDir", func() string {
		return Tpl(ToolDir)
	})
	Globals.Append("ToolPath", ToolPath)
}

func RunTools() {
	defer Handle()
	RunBinny()
	if FileExists(ToolPath("task")) {
		Cd(RootDir)
		NoErr(Exec(ToolPath("task"), ExecArgs(os.Args[1:]...)))
	}
}

func ToolPath(toolName string) string {
	toolPath := toolName
	switch runtime.GOOS {
	case "windows":
		toolPath += ".exe"
	}
	p := filepath.Join(Tpl(ToolDir), toolPath)
	Globals.Append("TOOL_"+strings.ToUpper(toolName), p)
	return p
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
