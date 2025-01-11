package tile

import (
	"fmt"

	"golang.org/x/exp/constraints"

	"github.com/hanselrd/domino/pkg/face"
)

type BasicTileFactory[T constraints.Integer] struct {
	ff       face.FaceFactory[T]
	numFaces uint
}

func NewBasicTileFactory[T constraints.Integer](ff face.FaceFactory[T], numFaces uint) BasicTileFactory[T] {
	return BasicTileFactory[T]{ff: ff, numFaces: numFaces}
}

func (btf BasicTileFactory[T]) CreateTile(vs ...T) (*Tile[T], error) {
	if len(vs) != int(btf.NumFaces()) {
		return nil, fmt.Errorf("%d values must be provided", btf.NumFaces())
	}
	return newTile(btf.ff, vs...)
}

func (btf BasicTileFactory[T]) FaceFactory() face.FaceFactory[T] {
	return btf.ff
}

func (btf BasicTileFactory[T]) NumFaces() uint {
	return btf.numFaces
}
