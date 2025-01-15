package domino

import (
	"fmt"
)

type TileFactory struct {
	ff       FaceFactory
	numFaces uint
}

func NewTileFactory(ff FaceFactory, numFaces uint) TileFactory {
	return TileFactory{ff: ff, numFaces: numFaces}
}

func (tf TileFactory) CreateTile(vs ...int) (*Tile, error) {
	if len(vs) != int(tf.NumFaces()) {
		return nil, fmt.Errorf("%d values must be provided", tf.NumFaces())
	}
	return newTile(tf.ff, vs...)
}

func (tf TileFactory) FaceFactory() FaceFactory {
	return tf.ff
}

func (tf TileFactory) NumFaces() uint {
	return tf.numFaces
}
