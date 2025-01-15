package domino

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

type Tile struct {
	faces []Face
}

func newTile(ff FaceFactory, vs ...int) (*Tile, error) {
	if len(vs) == 0 {
		return nil, fmt.Errorf("at least 1 value must be provided")
	}
	slices.Sort(vs)
	t2s := lo.Map(vs, func(v, _ int) lo.Tuple2[*Face, error] {
		return lo.T2(ff.CreateFace(v))
	})
	for _, t2 := range t2s {
		_, err := t2.Unpack()
		if err != nil {
			return nil, err
		}
	}
	return &Tile{faces: lo.Map(t2s, func(t2 lo.Tuple2[*Face, error], _ int) Face {
		return *t2.A
	})}, nil
}

func (t Tile) Faces() []Face {
	return t.faces
}

func (t Tile) IsMultiple() bool {
	return lo.EveryBy(t.faces, func(f Face) bool {
		return f == t.faces[0]
	})
}

func (t Tile) String() string {
	return strings.Join(lo.Map(t.faces, func(f Face, _ int) string {
		return f.String()
	}), ":")
}
