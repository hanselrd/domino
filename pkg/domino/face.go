package domino

import (
	"fmt"

	"github.com/kenshaw/baseconv"
	"github.com/samber/lo"
)

type Face int

func newFace(min, max, v int) (*Face, error) {
	if v < min || v > max {
		return nil, fmt.Errorf(
			"value %d is not within range [%d, %d]",
			v,
			min,
			max,
		)
	}
	f := Face(v)
	return &f, nil
}

func newUnsignedFace(max, v uint) (*Face, error) {
	return newFace(0, int(max), int(v))
}

func (f Face) NumPips() int {
	return int(f)
}

func (f Face) String() string {
	switch {
	case f.NumPips() >= 10:
		return lo.Must(baseconv.Encode64FromDec(fmt.Sprintf("%d", f.NumPips())))
	default:
		return fmt.Sprintf("%d", f.NumPips())
	}
}
