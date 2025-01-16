package domino

import (
	"slices"
	"sort"

	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/util/optionutil"
)

type TileSet struct {
	tiles []Tile
}

type TileSetOption func(*TileSet)

func WithShuffle() TileSetOption {
	return func(ts *TileSet) {
		ts.Shuffle()
	}
}

func NewTileSet(tf TileFactory, opts ...TileSetOption) (*TileSet, error) {
	ts := []Tile{}
	vs := lo.RangeFrom(
		tf.FaceFactory().MinValue(),
		int(tf.FaceFactory().MaxValue()-tf.FaceFactory().MinValue()+1),
	)
	vss := lo.RepeatBy(
		lo.Sum(
			lo.RepeatBy(int(tf.NumFaces()), func(i int) int {
				n := len(vs) - 1
				for range i {
					n *= len(vs)
				}
				return n
			}))+1,
		func(i int) []int {
			vz := []int{}
			for range tf.NumFaces() {
				vz = append(vz, i%len(vs))
				i /= len(vs)
			}
			slices.Sort(vz)
			return vz
		},
	)
	sort.Slice(vss, func(i, j int) bool {
		return slices.Compare(vss[i], vss[j]) < 0
	})
	vss = slices.CompactFunc(vss, func(a, b []int) bool {
		return slices.Compare(a, b) == 0
	})
	ts = lo.Map(vss, func(vs []int, _ int) Tile {
		return *lo.Must(tf.CreateTile(vs...))
	})
	return optionutil.Configure(&TileSet{tiles: ts}, opts), nil
}

func (ts TileSet) Tiles() []Tile {
	return ts.tiles
}

func (ts *TileSet) Shuffle() {
	lo.Shuffle(ts.tiles)
}

func (ts *TileSet) Draw(n int) []Tile {
	tz := make([]Tile, n)
	copy(tz, ts.tiles[:n])
	ts.tiles = slices.Delete(ts.tiles, 0, n)
	return tz
}

func (ts *TileSet) Return(tz ...Tile) {
	ts.tiles = slices.Concat(ts.tiles, tz)
}
