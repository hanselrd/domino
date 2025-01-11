package tile

import (
	"golang.org/x/exp/constraints"

	"github.com/hanselrd/domino/pkg/face"
)

type TileFactory[T constraints.Integer] interface {
	CreateTile(vs ...T) (*Tile[T], error)
	FaceFactory() face.FaceFactory[T]
	NumFaces() uint
}
