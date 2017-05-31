package class

import (
	"io"
	"encoding/binary"
	"errors"
	"go/constant"
)

const (
	ItemTop = iota
	ItemInteger
	ItemFloat
	ItemDouble
	ItemLong
	ItemNull
	ItemUninitializedThis
	ItemObject
	ItemUninitialized
)

type VerificationType interface {
	Tag() uint8
}

func NewVerificationType(io io.Reader, pool *ConstantPool) (VerificationType, error) {
	var tag uint8
	if err := binary.Read(io, binary.BigEndian, &tag); err != nil {
		return nil, err
	}
	switch tag {
	case ItemTop:
		return &TopVariable{}, nil
	case ItemInteger:
		return &IntegerVariable{}, nil
	case ItemFloat:
		return &FloatVariable{}, nil
	case ItemDouble:
		return &DoubleVariable{}, nil
	case ItemLong:
		return LongVariable{}, nil
	case ItemNull:
		return &NullVariable{}, nil
	case ItemUninitializedThis:
		return &UninitializedThisVariable{}, nil
	case ItemObject:
		return newObjectVariable(io, pool)
	case ItemUninitialized:
		return newUninitializedVariable(io, pool)
	}
	return nil, errors.New("unknow verification type tag")
}

type TopVariable struct {}
func (*TopVariable) Tag() uint8 {
	return ItemTop
}

type IntegerVariable struct {}
func (*IntegerVariable) Tag() uint8 {
	return ItemInteger
}

type FloatVariable struct {}
func (*FloatVariable) Tag() uint8 {
	return ItemFloat
}

type NullVariable struct {}
func (*NullVariable) Tag() uint8 {
	return ItemNull
}

type UninitializedThisVariable struct {}
func (*UninitializedThisVariable) Tag() uint8 {
	return ItemUninitializedThis
}

type ObjectVariable struct {
	Class *ConstClass
}

func newObjectVariable(io io.Reader, pool *ConstantPool) (*ObjectVariable, error) {
	var cpoolIndex uint16
	if err := binary.Read(io, binary.BigEndian, &cpoolIndex); err != nil {
		return nil, err
	}
	object := &ObjectVariable{}
	var err error
	if object.Class, err = pool.GetClass(cpoolIndex); err != nil {
		return nil, err
	}
	return object, nil
}

func (*ObjectVariable) Tag() uint8 {
	return ItemObject
}

type UninitializedVariable struct {
	Offset uint16
}

func newUninitializedVariable(io io.Reader, pool *ConstantPool) (*UninitializedVariable, error) {
	u := &UninitializedVariable{}
	if err := binary.Read(io, binary.BigEndian, &u.Offset); err != nil {
		return nil, err
	}
	return u, nil
}

func (*UninitializedVariable) Tag() uint8 {
	return ItemUninitialized
}

type LongVariable struct {}
func (*LongVariable) Tag() uint8 {
	return ItemLong
}

type DoubleVariable struct {}
func (*DoubleVariable) Tag() uint8 {
	return ItemDouble
}