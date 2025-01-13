package taskfile

type Task struct {
	Cmds      []Command    `yaml:"cmds,omitempty"`
	Deps      []Dependency `yaml:"deps,omitempty"`
	Desc      string       `yaml:"desc,omitempty"`
	Sources   []string     `yaml:"sources,omitempty"`
	Generates []string     `yaml:"generates,omitempty"`
	Requires  Requires     `yaml:"requires,omitempty"`
	Method    string       `yaml:"method,omitempty"`
}
