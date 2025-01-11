package face

import (
	"golang.org/x/exp/constraints"
)

type FaceFactory[T constraints.Integer] interface {
	CreateFace(v T) (*Face[T], error)
	MinValue() T
	MaxValue() T
}
