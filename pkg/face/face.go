package face

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Face[T constraints.Integer] struct {
	value T
}

func (f Face[T]) Value() T {
	return f.value
}

func (f Face[T]) String() string {
	return fmt.Sprintf("%d", f.Value())
}
