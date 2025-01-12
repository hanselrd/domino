package tile

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/constraints"

	"github.com/hanselrd/domino/pkg/face"
)

type Tile[T constraints.Integer] struct {
	faces []face.Face[T]
}

func newTile[T constraints.Integer](ff face.FaceFactory[T], vs ...T) (*Tile[T], error) {
	if len(vs) == 0 {
		return nil, fmt.Errorf("at least 1 value must be provided")
	}
	slices.Sort(vs)
	t2s := lo.Map(vs, func(v T, _ int) lo.Tuple2[*face.Face[T], error] {
		return lo.T2(ff.CreateFace(v))
	})
	for _, t2 := range t2s {
		_, err := t2.Unpack()
		if err != nil {
			return nil, err
		}
	}
	return &Tile[T]{faces: lo.Map(t2s, func(t2 lo.Tuple2[*face.Face[T], error], _ int) face.Face[T] {
		return *t2.A
	})}, nil
}

func (t Tile[T]) Faces() []face.Face[T] {
	return t.faces
}

func (t Tile[T]) IsMultiple() bool {
	return lo.EveryBy(t.faces, func(f face.Face[T]) bool {
		return f == t.faces[0]
	})
}

func (t Tile[T]) String() string {
	return strings.Join(lo.Map(t.faces, func(f face.Face[T], _ int) string {
		return f.String()
	}), ":")
}
