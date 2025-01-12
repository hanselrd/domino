package tileset

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"

	"github.com/hanselrd/domino/pkg/tile"
)

type TileSet struct {
	tiles []tile.Tile
}

func NewTileSet(tf tile.TileFactory) (*TileSet, error) {
	ts := []tile.Tile{}
	vs := lo.RangeFrom(tf.FaceFactory().MinValue(), int(tf.FaceFactory().MaxValue()-tf.FaceFactory().MinValue()+1))
	ns := lo.RangeFrom(0, int(lo.Must(strconv.ParseInt(strings.Repeat(strconv.FormatInt(int64(len(vs))-1, len(vs)), int(tf.NumFaces())), len(vs), 64)))+1)
	ss := lo.Map(ns, func(n, _ int) string {
		s := fmt.Sprintf("%0*s", tf.NumFaces(), strconv.FormatInt(int64(n), len(vs)))
		if len(s) != int(tf.NumFaces()) {
			panic(fmt.Sprintf("%s must have a length of %d", s, tf.NumFaces()))
		}
		return s
	})
	vss := lo.Map(ss, func(s string, _ int) []int {
		vs := lo.Map(strings.Split(s, ""), func(i string, _ int) int {
			return vs[lo.Must(strconv.ParseInt(i, len(vs), 64))]
		})
		slices.Sort(vs)
		return vs
	})
	sort.Slice(vss, func(i, j int) bool {
		return slices.Compare(vss[i], vss[j]) < 0
	})
	vss = slices.CompactFunc(vss, func(a, b []int) bool {
		return slices.Compare(a, b) == 0
	})
	ts = lo.Map(vss, func(vs []int, _ int) tile.Tile {
		return *lo.Must(tf.CreateTile(vs...))
	})
	return &TileSet{tiles: ts}, nil
}

func NewTileSetShuffled(tf tile.TileFactory) (*TileSet, error) {
	ts, err := NewTileSet(tf)
	if err != nil {
		return nil, err
	}
	ts.Shuffle()
	return ts, nil
}

func (ts TileSet) Tiles() []tile.Tile {
	return ts.tiles
}

func (ts *TileSet) Shuffle() {
	lo.Shuffle(ts.tiles)
}

func (ts *TileSet) Draw(n int) []tile.Tile {
	tz := make([]tile.Tile, n)
	copy(tz, ts.tiles[:n])
	ts.tiles = slices.Delete(ts.tiles, 0, n)
	return tz
}

func (ts *TileSet) Return(tz ...tile.Tile) {
	ts.tiles = slices.Concat(ts.tiles, tz)
}
