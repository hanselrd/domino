package face

import "fmt"

type FaceSigned = Face[int]

func newFaceSigned(min, max, v int) (*FaceSigned, error) {
	if v < min || v > max {
		return nil, fmt.Errorf("value %d is not within range [%d, %d]", v, min, max)
	}
	return &FaceSigned{value: v}, nil
}
