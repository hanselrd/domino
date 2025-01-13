package taskfile

type Taskfile struct {
	Version string              `yaml:"version,omitempty"`
	Vars    map[string]Variable `yaml:"vars,omitempty"`
	Tasks   map[string]Task     `yaml:"tasks,omitempty"`
}
