package class

import (
	"io"
	"encoding/binary"
	"fmt"
)

type Method struct {
	AccessFlags uint16
	Name *ConstUTF8
	Descriptor *ConstUTF8
	Attrs []Attr
}

func (method *Method) String() string {
	return fmt.Sprintf("Method -> %v", *method)
}

func NewMethod(io io.Reader, pool *ConstantPool) (*Method, error) {
	method := &Method{}
	if err := binary.Read(io, binary.BigEndian, &method.AccessFlags); err != nil {
		return nil, err
	}
	var nameIndex, descriptorIndex, attrCount uint16
	if err := binary.Read(io, binary.BigEndian, &nameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &descriptorIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &attrCount); err != nil {
		return nil, err
	}
	var err error
	if method.Name, err = pool.GetUTF8String(nameIndex); err != nil {
		return nil, err
	}
	if method.Descriptor, err = pool.GetUTF8String(descriptorIndex); err != nil {
		return nil, err
	}
	for i := 0; i < int(attrCount); i++ {
		attr, err := ReadAttr(io, pool)
		if err != nil {
			return nil, err
		}
		method.Attrs = append(method.Attrs, attr)
	}
	return method, nil
}