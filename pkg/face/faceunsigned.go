package face

import "fmt"

type FaceUnsigned = Face[uint]

func newFaceUnsigned(max, v uint) (*FaceUnsigned, error) {
	if v > max {
		return nil, fmt.Errorf("value %d is not within range [0, %d]", v, max)
	}
	return &FaceUnsigned{value: v}, nil
}
