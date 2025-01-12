package tileset_test

import (
	"testing"

	"github.com/hanselrd/domino/pkg/face"
	"github.com/hanselrd/domino/pkg/tile"
	"github.com/hanselrd/domino/pkg/tileset"
)

func TestTileSetT2F6(t *testing.T) {
	f6f := face.NewUnsignedFaceFactory(6)
	t2f := tile.NewTileFactory(f6f, 2)
	ts, err := tileset.NewTileSetShuffled(t2f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 28 {
		t.Error(ln)
	}
}

func TestTileSetT2F9(t *testing.T) {
	f9f := face.NewUnsignedFaceFactory(9)
	t2f := tile.NewTileFactory(f9f, 2)
	ts, err := tileset.NewTileSetShuffled(t2f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 55 {
		t.Error(ln)
	}
}

func TestTileSetT2F12(t *testing.T) {
	f12f := face.NewUnsignedFaceFactory(12)
	t2f := tile.NewTileFactory(f12f, 2)
	ts, err := tileset.NewTileSetShuffled(t2f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 91 {
		t.Error(ln)
	}
}

func TestTileSetT2F15(t *testing.T) {
	f15f := face.NewUnsignedFaceFactory(15)
	t2f := tile.NewTileFactory(f15f, 2)
	ts, err := tileset.NewTileSetShuffled(t2f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 136 {
		t.Error(ln)
	}
}

func TestTileSetT3F6(t *testing.T) {
	f6f := face.NewUnsignedFaceFactory(6)
	t3f := tile.NewTileFactory(f6f, 3)
	ts, err := tileset.NewTileSet(t3f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 84 {
		t.Error(ln)
	}
}

func TestTileSetT3F9(t *testing.T) {
	f9f := face.NewUnsignedFaceFactory(9)
	t3f := tile.NewTileFactory(f9f, 3)
	ts, err := tileset.NewTileSet(t3f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 220 {
		t.Error(ln)
	}
}

func TestTileSetT3F12(t *testing.T) {
	f12f := face.NewUnsignedFaceFactory(12)
	t3f := tile.NewTileFactory(f12f, 3)
	ts, err := tileset.NewTileSet(t3f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 455 {
		t.Error(ln)
	}
}

func TestTileSetT3F15(t *testing.T) {
	f15f := face.NewUnsignedFaceFactory(15)
	t3f := tile.NewTileFactory(f15f, 3)
	ts, err := tileset.NewTileSet(t3f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	if ln := len(ts.Tiles()); ln != 816 {
		t.Error(ln)
	}
}

func TestTileSetT2F6DrawReturn(t *testing.T) {
	f6f := face.NewUnsignedFaceFactory(6)
	t2f := tile.NewTileFactory(f6f, 2)
	ts, err := tileset.NewTileSetShuffled(t2f)
	if err != nil {
		t.Error(err)
	}
	t.Logf("created %d tile(s)", len(ts.Tiles()))
	preln := len(ts.Tiles())
	tz := ts.Draw(7)
	t.Logf("drew %d tile(s): %s", len(tz), tz)
	if ln := len(ts.Tiles()); ln != preln-len(tz) {
		t.Error(ln)
	}
	if ln := len(tz); ln != len(tz) {
		t.Error(ln)
	}
	ts.Return(tz...)
	t.Logf("returned %d tile(s): %s", len(tz), tz)
	if ln := len(ts.Tiles()); ln != preln {
		t.Error(ln)
	}
}
