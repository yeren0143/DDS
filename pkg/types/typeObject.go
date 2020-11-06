package types

type CompleteTypeObject struct {
}

func NewCompleteTypeObject() *CompleteTypeObject {
	return &CompleteTypeObject{}
}

type MinimalTypeObject struct {
}

func NewMinimalTypeObject() *MinimalTypeObject {
	return &MinimalTypeObject{}
}

type TypeObject_t struct {
	Md          uint8
	CompleteObj *CompleteTypeObject
	MinimalObj  *MinimalTypeObject
}

func NewTypeObject() *TypeObject_t {
	return &TypeObject_t{
		Md:          0x00,
		CompleteObj: NewCompleteTypeObject(),
		MinimalObj:  NewMinimalTypeObject(),
	}
}
