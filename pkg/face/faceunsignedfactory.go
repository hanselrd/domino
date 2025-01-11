package face

type FaceUnsignedFactory = FaceFactory[uint]

type faceUnsignedFactory struct {
	max uint
}

func NewFaceUnsignedFactory(max uint) FaceUnsignedFactory {
	return faceUnsignedFactory{max: max}
}

func (fuf faceUnsignedFactory) CreateFace(v uint) (*FaceUnsigned, error) {
	return newFaceUnsigned(fuf.MaxValue(), v)
}

func (fuf faceUnsignedFactory) MinValue() uint {
	return 0
}

func (fuf faceUnsignedFactory) MaxValue() uint {
	return fuf.max
}
