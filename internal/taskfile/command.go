package taskfile

type Command interface{}

type CommandString string

type CommandStruct struct {
	Task string `yaml:"task,omitempty"`
}
