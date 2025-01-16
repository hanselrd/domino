package taskfile

type Taskfile struct {
	Version string              `yaml:"version,omitempty"`
	Method  string              `yaml:"method,omitempty"`
	Vars    map[string]Variable `yaml:"vars,omitempty"`
	Tasks   map[string]Task     `yaml:"tasks,omitempty"`
	Set     []string            `yaml:"set,omitempty"`
	Shopt   []string            `yaml:"shopt,omitempty"`
}
