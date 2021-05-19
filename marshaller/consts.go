package marshaller

const (
	TypeIdBool uint16 = iota
	TypeIdInt
	TypeIdFloat32
	TypeIdFloat64
	TypeIdString
)

var (
	TypeIdSlice = []byte{0xFF, 0xFF}
	TypeIdMap   = []byte{0xFF, 0xFE}
)
