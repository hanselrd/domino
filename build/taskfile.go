//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/hanselrd/domino/internal/taskfile"
	"github.com/hanselrd/domino/internal/util/maputil"
)

var (
	bins      = []string{"domino"}
	platforms = []string{
		"windows/amd64",
		"linux/amd64",
		"darwin/amd64",
		"darwin/arm64",
	}
	builds  = []string{"debug", "release"}
	sources = []string{"cmd/**/*.go", "internal/**/*.go", "pkg/**/*.go"}
	gcflags = map[string]string{
		"debug":   "all=-N -l",
		"release": "all=-l -B -C",
	}
	ldflags = map[string]string{
		"debug":   "",
		"release": "-s -w",
	}
	buildMetadataVars = lo.MapKeys(map[string]taskfile.Variable{
		"VERSION": "0.0.1-alpha.1",
		"TIME": taskfile.VariableDynamic{
			Sh: "date --utc \"+%Y-%m-%dT%H:%M:%SZ\"",
		},
		"HASH": taskfile.VariableDynamic{
			Sh: "git rev-parse HEAD",
		},
		"SHORT_HASH": taskfile.VariableDynamic{
			Sh: "git rev-parse --short=7 HEAD",
		},
	}, func(_ taskfile.Variable, k string) string {
		return fmt.Sprintf("BUILD_%s", k)
	})
	buildMetadataVarNames = lo.Must(maputil.SortedKeys(buildMetadataVars))
	buildMetadataLdflags  = strings.Join(
		lo.Map(buildMetadataVarNames, func(n string, _ int) string {
			return fmt.Sprintf(
				"-X 'github.com/hanselrd/domino/internal/build.%s={{.%s}}'",
				lo.PascalCase(strings.TrimPrefix(n, "BUILD_")),
				n,
			)
		}),
		" ",
	)
)

func osArch(platform string) (os, arch string) {
	s := strings.Split(platform, "/")
	if len(s) != 2 {
		panic(s)
	}
	os, arch = s[0], s[1]
	return
}

var tf = taskfile.Taskfile{
	Version: "3",
	Vars:    buildMetadataVars,
	Tasks: func() map[string]taskfile.Task {
		ts0 := map[string]taskfile.Task{
			"default": {
				Cmds: lo.Map(
					[]string{"format", "bootstrap", "build"},
					func(t string, _ int) taskfile.Command {
						return taskfile.CommandStruct{
							Task: t,
						}
					},
				),
			},
			"bootstrap": {
				Cmds: []taskfile.Command{
					"go run build/taskfile.go",
				},
				Sources:   []string{"build/taskfile.go"},
				Generates: []string{"Taskfile.yml"},
				Method:    "timestamp",
			},
			"build": {
				Deps: lo.Map(builds, func(b string, _ int) taskfile.Dependency {
					return fmt.Sprintf("build-%s", b)
				}),
			},
			"format": {
				Cmds: []taskfile.Command{
					"goimports -w -local \"github.com/hanselrd/domino\" .",
					"gofumpt -w -extra .",
					"golines -w -m 100 **/*.go",
				},
				Sources: []string{"**/*.go"},
			},
			"update": {
				Cmds: []taskfile.Command{
					"go get -u ./...",
					"go mod tidy",
					"go get gopkg.in/yaml.v3",
				},
			},
			"test": {
				Cmds: []taskfile.Command{
					fmt.Sprintf(
						"go test -gcflags=\"%s\" -ldflags=\"%s\" -v ./...",
						gcflags["debug"],
						ldflags["debug"],
					),
				},
			},
			"clean": {
				Cmds: []taskfile.Command{"rm -rf bin"},
			},
		}
		ts1 := map[string]taskfile.Task{}
		for _, b := range bins {
			for _, p := range platforms {
				os, arch := osArch(p)
				for _, bb := range builds {
					ts1[fmt.Sprintf("build-%s-%s-%s-%s", b, os, arch, bb)] = taskfile.Task{
						Cmds: []taskfile.Command{
							fmt.Sprintf("mkdir -p bin/%s", bb),
							fmt.Sprintf(
								"GOOS=%[2]s GOARCH=%[3]s go build -gcflags=\"%[5]s\" -ldflags=\"%[6]s\" -o bin/%[4]s/%[1]s_%[2]s_%[3]s%[7]s ./cmd/%[1]s",
								b,
								os,
								arch,
								bb,
								gcflags[bb],
								strings.TrimSpace(
									strings.Join([]string{ldflags[bb], buildMetadataLdflags}, " "),
								),
								lo.Ternary(os == "windows", ".exe", ""),
							),
						},
						Sources: sources,
						Generates: []string{
							fmt.Sprintf(
								"bin/%s/%s_%s_%s%s",
								bb,
								b,
								os,
								arch,
								lo.Ternary(os == "windows", ".exe", ""),
							),
						},
						Requires: taskfile.Requires{
							Vars: buildMetadataVarNames,
						},
					}
				}
			}
			for _, bb := range builds {
				ts1[fmt.Sprintf("build-%s-%s", b, bb)] = taskfile.Task{
					Cmds: []taskfile.Command{
						fmt.Sprintf("mkdir -p bin/%s", bb),
						fmt.Sprintf(
							"go build -gcflags=\"%[3]s\" -ldflags=\"%[4]s\" -o bin/%[2]s/%[1]s%[5]s ./cmd/%[1]s",
							b,
							bb,
							gcflags[bb],
							strings.TrimSpace(
								strings.Join([]string{ldflags[bb], buildMetadataLdflags}, " "),
							),
							lo.Ternary(runtime.GOOS == "windows", ".exe", ""),
						),
					},
					Sources: sources,
					Generates: []string{
						fmt.Sprintf(
							"bin/%s/%s%s",
							bb,
							b,
							lo.Ternary(runtime.GOOS == "windows", ".exe", ""),
						),
					},
					Requires: taskfile.Requires{
						Vars: buildMetadataVarNames,
					},
				}
			}
		}
		ts2 := map[string]taskfile.Task{}
		for _, b := range builds {
			ks := lo.Filter(lo.Must(maputil.SortedKeys(ts1)), func(k string, _ int) bool {
				return strings.HasSuffix(k, fmt.Sprintf("-%s", b))
			})
			ts2[fmt.Sprintf("build-%s", b)] = taskfile.Task{
				Deps: lo.Map(ks, func(k string, _ int) taskfile.Dependency {
					return k
				}),
			}
		}
		ts := lo.Assign(ts0, ts1, ts2)
		return ts
	}(),
}

func main() {
	b := bytes.NewBufferString("")
	enc := yaml.NewEncoder(b)
	enc.SetIndent(2)
	lo.Must0(enc.Encode(tf))
	lo.Must0(os.WriteFile("Taskfile.yml", []byte(b.String()), 0o644))
}
