//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type taskfile struct {
	Version string                  `yaml:"version,omitempty"`
	Tasks   map[string]taskfileTask `yaml:"tasks,omitempty"`
}

type taskfileTask struct {
	Name string   `yaml:"name,omitempty"`
	Deps []string `yaml:"deps,omitempty"`
	Cmds []string `yaml:"cmds,omitempty"`
}

var bins = []string{"domino"}
var platforms = []string{"windows/amd64", "linux/amd64", "darwin/amd64", "darwin/arm64"}
var builds = []string{"debug", "release"}
var gcflags = map[string]string{
	"debug":   "all=-N -l",
	"release": "",
}
var ldflags = map[string]string{
	"debug":   "",
	"release": "-s -w",
}

func osArch(platform string) (os, arch string) {
	s := strings.Split(platform, "/")
	if len(s) != 2 {
		panic(s)
	}
	os, arch = s[0], s[1]
	return
}

var tf = taskfile{
	Version: "3",
	Tasks: func() (ts map[string]taskfileTask) {
		ts = map[string]taskfileTask{
			"bootstrap": taskfileTask{
				Deps: []string{"update"},
				Cmds: []string{
					"go run build/taskfile.go",
				}},
			"build": taskfileTask{
				Deps: lo.Map(builds, func(b string, _ int) string {
					return fmt.Sprintf("build-%s", b)
				}),
			},
			"format": taskfileTask{
				Cmds: []string{
					"goimports -w -local \"github.com/hanselrd/domino\" .",
				}},
			"update": taskfileTask{
				Cmds: []string{
					"go get -u ./...",
					"go mod tidy",
					"go get gopkg.in/yaml.v3",
				}},
			"test": taskfileTask{
				Cmds: []string{
					fmt.Sprintf("go test -gcflags=\"%s\" -ldflags=\"%s\" -v ./...", gcflags["debug"], ldflags["debug"]),
				}},
			"clean": taskfileTask{
				Cmds: []string{"rm -rf bin"},
			}}
		tz := map[string]taskfileTask{}
		for _, b := range bins {
			for _, p := range platforms {
				os, arch := osArch(p)
				for _, bb := range builds {
					tz[fmt.Sprintf("build-%s-%s-%s-%s", b, os, arch, bb)] = taskfileTask{
						Cmds: []string{
							fmt.Sprintf("mkdir -p bin/%s", bb),
							fmt.Sprintf("GOOS=%[2]s GOARCH=%[3]s go build -gcflags=\"%[5]s\" -ldflags=\"%[6]s\" -o bin/%[4]s/%[1]s_%[2]s_%[3]s%[7]s ./cmd/%[1]s", b, os, arch, bb, gcflags[bb], ldflags[bb], lo.Ternary(os == "windows", ".exe", "")),
						},
					}
				}
			}
		}
		for k, v := range tz {
			ts[k] = v
		}
		for _, b := range builds {
			ts[fmt.Sprintf("build-%s", b)] = taskfileTask{
				Deps: lo.Keys(tz),
			}
		}
		for _, v := range ts {
			slices.Sort(v.Deps)
		}
		return
	}(),
}

func main() {
	b := bytes.NewBufferString("")
	enc := yaml.NewEncoder(b)
	enc.SetIndent(2)
	lo.Must0(enc.Encode(tf))
	lo.Must0(os.WriteFile("Taskfile.yml", []byte(b.String()), 0o644))
}
