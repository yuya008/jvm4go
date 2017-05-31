package class

type Field struct {
	AccessFlags uint16
	Name *ConstUTF8
	Descriptor *ConstUTF8
	Attrs []Attr
}
