package build

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

func RunBinny() {
	binnyPath := ToolPath("binny")
	if !FileExists(binnyPath) {
		installBinny(binnyPath)
	}
	Cd(RootDir) // could be: FindFile(".binny.yaml")
	Run(binnyPath, "install", "-v")
}

func installBinny(binnyPath string) {
	version := findBinnyVersion()

	err := downloadPrebuiltBinary(binnyPath, downloadSpec{
		url: "https://github.com/anchore/binny/releases/download/v{{.version}}/binny_{{.version}}_{{.os}}_{{.arch}}.{{.ext}}",
		args: map[string]string{
			"ext":     "tar.gz",
			"version": version,
		},
		platform: map[string]map[string]string{
			"windows": {
				"ext": "zip",
			},
		},
	})

	if err != nil {
		LogErr(err)

		BuildFromGoSource(
			binnyPath,
			"github.com/anchore/binny",
			"cmd/binny",
			version,
			Ldflags("-w",
				"-s",
				"-extldflags '-static'",
				"-X main.version="+GoDepVersion("github.com/anchore/binny")))
	}
}

func findBinnyVersion() string {
	binnyConfig := FindFile(".binny.yaml")
	if binnyConfig != "" {
		cfg := map[string]any{}
		f, err := os.Open(binnyConfig)
		if err == nil {
			d := yaml.NewDecoder(f)
			err = d.Decode(&cfg)
			if err == nil {
				tools := cfg["tools"]
				if tools, ok := tools.([]any); ok {
					for _, tool := range tools {
						if m, ok := tool.(map[string]any); ok && m["name"] == "binny" {
							v := m["version"]
							if v, ok := v.(string); ok {
								return v
							}
							if v, ok := v.(map[string]any); ok {
								if v, ok := v["want"].(string); ok {
									return regexp.MustCompile("^v").ReplaceAllString(v, "")
								}
							}
						}
					}
				}
			}
		}
	}
	return "0.8.0"
}

func downloadPrebuiltBinary(toolPath string, spec downloadSpec) error {
	tplArgs := spec.currentArgs()
	url := Tpl(spec.url, tplArgs)
	contents, code, status := Fetch(url)
	if code > 300 || len(contents) == 0 {
		return fmt.Errorf("error downloading %v: http %v %v", url, code, status)
	}
	contents = getArchiveFileContents(contents, filepath.Base(toolPath))
	if contents == nil {
		return fmt.Errorf("unable to read archive from downloading %v: http %v %v", url, code, status)
	}
	dir := filepath.Dir(toolPath)
	if !FileExists(dir) {
		NoErr(os.MkdirAll(dir, 0700|os.ModeDir))
	}
	return os.WriteFile(toolPath, contents, 0500) // read + execute permissions
}

func getArchiveFileContents(archive []byte, file string) []byte {
	var errs []error

	contents, err := getZipArchiveFileContents(archive, file)
	if err == nil && len(contents) > 0 {
		return contents
	}
	errs = append(errs, err)

	contents, err = getTarGzArchiveFileContents(archive, file)
	if err == nil && len(contents) > 0 {
		return contents
	}
	errs = append(errs, err)

	Throw(fmt.Errorf("unable to read archive after attempting readers: %w", errors.Join(errs...)))
	return nil
}

func getZipArchiveFileContents(archive []byte, file string) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}
	f, err := zipReader.Open(file)
	if err != nil {
		return nil, err
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func getTarGzArchiveFileContents(archive []byte, file string) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(archive))
	if err == nil && gzipReader != nil {
		t := tar.NewReader(gzipReader)
		for {
			hdr, err := t.Next()
			if err != nil {
				return nil, err
			}
			if hdr.Name == file {
				const GB = 1024 * 1024 * 1024
				if hdr.Size > 2*GB {
					return nil, fmt.Errorf("refusing to extract file %v larger than 2 GB, declared size: %v", file, hdr.Size)
				}
				return io.ReadAll(t)
			}
		}
	}
	return nil, fmt.Errorf("file not found: %v", file)
}

func Ldflags(flags ...string) ExecOpt {
	return func(cmd *exec.Cmd) {
		for i, arg := range cmd.Args {
			// append to existing ldflags arg
			if arg == "-ldflags" {
				if i+1 >= len(cmd.Args) {
					cmd.Args = append(cmd.Args, "")
				} else {
					cmd.Args[i+1] += " "
				}
				cmd.Args[i+1] += strings.Join(flags, " ")
				return
			}
		}
		cmd.Args = append(cmd.Args, "-ldflags", strings.Join(flags, " "))
	}
}

func BuildFromGoSource(file, module, entrypoint, version string, opts ...ExecOpt) {
	Log("Building: %s", module)
	InGitClone("https://"+module, version, func() {
		NoErr(Exec("go", ExecArgs("build"), ExecOpts(opts...), ExecArgs("-o", file, "./"+entrypoint), ExecStd()))
	})
}

type downloadSpec struct {
	url      string
	args     map[string]string
	platform map[string]map[string]string
}

func (d downloadSpec) currentArgs() map[string]any {
	out := map[string]any{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
	for k, v := range d.args {
		out[k] = v
	}
	if d.platform != nil {
		for k, v := range d.platform[runtime.GOOS] {
			out[k] = v
		}
	}
	return out
}
