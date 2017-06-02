package class

type Field struct {
	AccessFlags uint16
	Name *ConstUTF8
	Descriptor *ConstUTF8
	Attrs []Attr
}

const (
	FieldAccPublic = 0x0001
	FieldAccPrivate = 0x0002
	FieldAccProtected = 0x0004
	FieldAccStatic = 0x0008
	FieldAccFinal = 0x0010
	FieldAccVolatile = 0x0040
	FieldAccTransient = 0x0080
	FieldAccSynthetic = 0x1000
	FieldAccEnum = 0x4000
)