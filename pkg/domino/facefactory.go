package domino

type FaceFactory struct {
	min, max int
}

func NewFaceFactory(min, max int) FaceFactory {
	return FaceFactory{min: min, max: max}
}

func NewUnsignedFaceFactory(max uint) FaceFactory {
	return NewFaceFactory(0, int(max))
}

func (ff FaceFactory) CreateFace(v int) (*Face, error) {
	return newFace(ff.min, ff.max, v)
}

func (ff FaceFactory) MinValue() int {
	return ff.min
}

func (ff FaceFactory) MaxValue() int {
	return ff.max
}
