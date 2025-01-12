package tile

import (
	"fmt"

	"github.com/hanselrd/domino/pkg/face"
)

type TileFactory struct {
	ff       face.FaceFactory
	numFaces uint
}

func NewTileFactory(ff face.FaceFactory, numFaces uint) TileFactory {
	return TileFactory{ff: ff, numFaces: numFaces}
}

func (tf TileFactory) CreateTile(vs ...int) (*Tile, error) {
	if len(vs) != int(tf.NumFaces()) {
		return nil, fmt.Errorf("%d values must be provided", tf.NumFaces())
	}
	return newTile(tf.ff, vs...)
}

func (tf TileFactory) FaceFactory() face.FaceFactory {
	return tf.ff
}

func (tf TileFactory) NumFaces() uint {
	return tf.numFaces
}
