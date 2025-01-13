package taskfile

type Variable interface{}

type (
	VariableString string
	VariableStatic = VariableString
)

type VariableStruct struct {
	Sh string `yaml:"sh"`
}
type VariableDynamic = VariableStruct
