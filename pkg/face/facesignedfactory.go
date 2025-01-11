package face

type FaceSignedFactory = FaceFactory[int]

type faceSignedFactory struct {
	min, max int
}

func NewFaceSignedFactory(min, max int) FaceSignedFactory {
	return faceSignedFactory{min: min, max: max}
}

func (fsf faceSignedFactory) CreateFace(v int) (*FaceSigned, error) {
	return newFaceSigned(fsf.MinValue(), fsf.MaxValue(), v)
}

func (fsf faceSignedFactory) MinValue() int {
	return fsf.min
}

func (fsf faceSignedFactory) MaxValue() int {
	return fsf.max
}
