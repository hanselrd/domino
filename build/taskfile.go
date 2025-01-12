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
	Vars    map[string]taskfileVar  `yaml:"vars,omitempty"`
	Tasks   map[string]taskfileTask `yaml:"tasks,omitempty"`
}

type (
	taskfileVar        interface{}
	taskfileVarStatic  string
	taskfileVarDynamic struct {
		Sh string `yaml:"sh"`
	}
)

type taskfileTask struct {
	Deps     []string         `yaml:"deps,omitempty"`
	Cmds     []string         `yaml:"cmds,omitempty"`
	Label    string           `yaml:"label,omitempty"`
	Desc     string           `yaml:"desc,omitempty"`
	Requires taskfileRequires `yaml:"requires,omitempty"`
}

type taskfileRequires struct {
	Vars []string `yaml:"vars"`
}

var (
	bins      = []string{"domino"}
	platforms = []string{
		"windows/amd64",
		"linux/amd64",
		"darwin/amd64",
		"darwin/arm64",
	}
	builds = []string{"debug", "release"}
)

var (
	gcflags = map[string]string{
		"debug":   "all=-N -l",
		"release": "all=-l -B -C",
	}
	ldflags = map[string]string{
		"debug":   "",
		"release": "-s -w",
	}
)

func osArch(platform string) (os, arch string) {
	s := strings.Split(platform, "/")
	if len(s) != 2 {
		panic(s)
	}
	os, arch = s[0], s[1]
	return
}

var (
	buildMetadataVars = map[string]taskfileVar{
		"BUILD_VERSION": "0.0.1-alpha.1",
		"BUILD_TIME": taskfileVarDynamic{
			Sh: "date --utc",
		},
		"BUILD_HASH": taskfileVarDynamic{
			Sh: "git rev-parse HEAD",
		},
		"BUILD_SHORT_HASH": taskfileVarDynamic{
			Sh: "git rev-parse --short=6 HEAD",
		},
	}
	buildMetadataLdflags = func() string {
		pkg := "github.com/hanselrd/domino/internal/build"
		return strings.Join([]string{
			fmt.Sprintf("-X '%s.Version={{.BUILD_VERSION}}'", pkg),
			fmt.Sprintf("-X '%s.Time={{.BUILD_TIME}}'", pkg),
			fmt.Sprintf("-X '%s.Hash={{.BUILD_HASH}}'", pkg),
			fmt.Sprintf("-X '%s.ShortHash={{.BUILD_SHORT_HASH}}'", pkg),
		}, " ")
	}()
)

var tf = taskfile{
	Version: "3",
	Vars:    buildMetadataVars,
	Tasks: func() (ts map[string]taskfileTask) {
		ts = map[string]taskfileTask{
			"bootstrap": {
				Deps: []string{"update"},
				Cmds: []string{
					"go run build/taskfile.go",
				},
			},
			"default": {
				Deps: []string{"build"},
			},
			"build": {
				Deps: lo.Map(builds, func(b string, _ int) string {
					return fmt.Sprintf("build-%s", b)
				}),
			},
			"format": {
				Cmds: []string{
					"goimports -w -local \"github.com/hanselrd/domino\" .",
					"gofumpt -w -extra .",
				},
			},
			"update": {
				Cmds: []string{
					"go get -u ./...",
					"go mod tidy",
					"go get gopkg.in/yaml.v3",
				},
			},
			"test": {
				Cmds: []string{
					fmt.Sprintf("go test -gcflags=\"%s\" -ldflags=\"%s\" -v ./...", gcflags["debug"], ldflags["debug"]),
				},
			},
			"clean": {
				Cmds: []string{"rm -rf bin"},
			},
		}
		tz := map[string]taskfileTask{}
		for _, b := range bins {
			for _, p := range platforms {
				os, arch := osArch(p)
				for _, bb := range builds {
					tz[fmt.Sprintf("build-%s-%s-%s-%s", b, os, arch, bb)] = taskfileTask{
						Cmds: []string{
							fmt.Sprintf("mkdir -p bin/%s", bb),
							fmt.Sprintf("GOOS=%[2]s GOARCH=%[3]s go build -gcflags=\"%[5]s\" -ldflags=\"%[6]s\" -o bin/%[4]s/%[1]s_%[2]s_%[3]s%[7]s ./cmd/%[1]s", b, os, arch, bb, gcflags[bb], strings.TrimSpace(strings.Join([]string{ldflags[bb], buildMetadataLdflags}, " ")), lo.Ternary(os == "windows", ".exe", "")),
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
		for _, k := range lo.Keys(ts) {
			slices.Sort(ts[k].Deps)
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
