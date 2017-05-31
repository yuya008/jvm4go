package class

import (
	"io"
	"errors"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
	"math"
	"reflect"
)

const (
	UTF8               = 1
	Integer            = 3
	Float              = 4
	Long               = 5
	Double             = 6
	Class              = 7
	String             = 8
	FieldRef           = 9
	MethodRef          = 10
	InterfaceMethodRef = 11
	NameAndType        = 12
	MethodHandle       = 15
	MethodType         = 16
	InvokeDynamic      = 18
)

const (
	_ = iota
	RefGetField
	RefGetStatic
	RefPutField
	RefPutStatic
	RefInvokeVirtual
	RefInvokeStatic
	RefInvokeSpecial
	RefNewInvokeSpecial
	RefInvokeInterface
)

var (
	ResolvingError = errors.New("resolving fail")
)

type Constant interface {
	Tag() int
	Resolving(*ConstantPool) error
}

type ConstantPool struct {
	pool []Constant
}

func NewConstantPool(reader io.Reader) (*ConstantPool, error) {
	var constantCount uint16
	if err := binary.Read(reader, binary.BigEndian, &constantCount); err != nil {
		return nil, err
	}
	constPool := &ConstantPool{pool: make([]Constant, constantCount)}
	var tag uint8
	for i := 1; i < int(constantCount); i++ {
		if err := binary.Read(reader, binary.BigEndian, &tag); err != nil {
			return nil, err
		}
		constant, err := NewConstant(tag, reader)
		if err != nil {
			return nil, err
		}
		constPool.pool[i] = constant
		if tag == Long || tag == Double {
			i++
		}
	}
	// Resolving
	for i := 1; i < int(constantCount); i++ {
		if c := constPool.pool[i]; c != nil {
			if err := c.Resolving(constPool); err != nil {
				return nil, err
			}
		}
	}
	return constPool, nil
}

func (cp *ConstantPool) Get(i uint16) (Constant, error) {
	if int(i) < 0 || int(i) > len(cp.pool) - 1 {
		return nil, fmt.Errorf("index is %d constant not found", i)
	}
	constant := cp.pool[i]
	if constant == nil {
		return nil, fmt.Errorf("index is %d constant not found", i)
	}
	return constant, nil
}

func (cp *ConstantPool) GetUTF8String(i uint16) (*ConstUTF8, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstUTF8); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstUTF8) failed")
}

func (cp *ConstantPool) GetFloat(i uint16) (*ConstFloat, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstFloat); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstFloat) failed")
}

func (cp *ConstantPool) GetInteger(i uint16) (*ConstInteger, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstInteger); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstInteger) failed")
}

func (cp *ConstantPool) GetLong(i uint16) (*ConstLong, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstLong); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstLong) failed")
}

func (cp *ConstantPool) GetDouble(i uint16) (*ConstDouble, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstDouble); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstDouble) failed")
}

func (cp *ConstantPool) GetFieldRef(i uint16) (*ConstFieldRef, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstFieldRef); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstFieldRef) failed")
}

func (cp *ConstantPool) GetString(i uint16) (*ConstString, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstString); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstString) failed")
}

func (cp *ConstantPool) GetClass(i uint16) (*ConstClass, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstClass); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstClass) failed")
}

func (cp *ConstantPool) GetMethodRef(i uint16) (*ConstMethodRef, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstMethodRef); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstMethodRef) failed")
}

func (cp *ConstantPool) GetInterfaceMethodRef(i uint16) (*ConstInterfaceMethodRef, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstInterfaceMethodRef); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstInterfaceMethodRef) failed")
}

func (cp *ConstantPool) GetNameAndType(i uint16) (*ConstNameAndType, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstNameAndType); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstNameAndType) failed")
}

func (cp *ConstantPool) GetMethodHandle(i uint16) (*ConstMethodHandle, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstMethodHandle); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstMethodHandle) failed")
}

func (cp *ConstantPool) GetMethodType(i uint16) (*ConstMethodType, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstMethodType); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstMethodType) failed")
}

func (cp *ConstantPool) GetInvokeDynamic(i uint16) (*ConstInvokeDynamic, error) {
	constant, err := cp.Get(i)
	if err != nil {
		return nil, err
	}
	if val, ok := constant.(*ConstInvokeDynamic); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstInvokeDynamic) failed")
}

func (cp *ConstantPool) Length() int {
	return len(cp.pool)
}

func NewConstant(tag uint8, reader io.Reader) (Constant, error) {
	switch tag {
	case UTF8:
		return NewConstUTF8(reader)
	case Integer:
		return NewConstInteger(reader)
	case Float:
		return NewConstFloat(reader)
	case Long:
		return NewConstLong(reader)
	case Double:
		return NewConstDouble(reader)
	case Class:
		return NewConstClass(reader)
	case String:
		return NewConstString(reader)
	case FieldRef:
		return NewConstFieldRef(reader)
	case MethodRef:
		return NewConstMethodRef(reader)
	case InterfaceMethodRef:
		return NewConstInterfaceMethodRef(reader)
	case NameAndType:
		return NewConstNameAndType(reader)
	case MethodHandle:
		return NewConstMethodHandle(reader)
	case MethodType:
		return NewConstMethodType(reader)
	case InvokeDynamic:
		return NewConstInvokeDynamic(reader)
	}
	return nil, errors.New("unknow constant tag")
}

type ConstUTF8 struct {
	s string
}

func NewConstUTF8(io io.Reader) (*ConstUTF8, error) {
	cu := &ConstUTF8{}
	var length uint16
	if err := binary.Read(io, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	if length <= 0 {
		return nil, ClassFileFormatError
	}
	buf := make([]byte, length)
	if err := binary.Read(io, binary.BigEndian, buf); err != nil {
		return nil, err
	}
	var err error
	if cu.s, err = decodeMUTF8(buf); err != nil {
		return nil, err
	}
	return cu, nil
}

func (c *ConstUTF8) Tag() int {
	return UTF8
}

func (c *ConstUTF8) Resolving(pool *ConstantPool) error {
	return nil
}

func (c *ConstUTF8) String() string {
	return c.s
}

type ConstInteger struct {
	Val int32
}

func NewConstInteger(io io.Reader) (*ConstInteger, error) {
	ci := &ConstInteger{}
	if err := binary.Read(io, binary.BigEndian, &ci.Val); err != nil {
		return nil, err
	}
	return ci, nil
}

func (*ConstInteger) Tag() int {
	return Integer
}

func (*ConstInteger) Resolving(pool *ConstantPool) error {
	return nil
}

type ConstFloat struct {
	Val float32
}

func NewConstFloat(io io.Reader) (*ConstFloat, error) {
	cf := &ConstFloat{}
	var v uint32
	if err := binary.Read(io, binary.BigEndian, &v); err != nil {
		return nil, err
	}
	cf.Val = math.Float32frombits(v)
	return cf, nil
}

func (*ConstFloat) Tag() int {
	return Float
}

func (cf *ConstFloat) Resolving(pool *ConstantPool) error {
	return nil
}

type ConstLong struct {
	Val int64
}

func NewConstLong(io io.Reader) (*ConstLong, error) {
	cl := &ConstLong{}
	if err := binary.Read(io, binary.BigEndian, &cl.Val); err != nil {
		return nil, err
	}
	return cl, nil
}

func (*ConstLong) Tag() int {
	return Long
}

func (cl *ConstLong) Resolving(pool *ConstantPool) error {
	return nil
}

type ConstDouble struct {
	Val float64
}

func NewConstDouble(io io.Reader) (*ConstDouble, error) {
	cd := &ConstDouble{}
	var v uint64
	if err := binary.Read(io, binary.BigEndian, &v); err != nil {
		return nil, err
	}
	cd.Val = math.Float64frombits(v)
	return cd, nil
}

func (*ConstDouble) Tag() int {
	return Double
}

func (cl *ConstDouble) Resolving(pool *ConstantPool) error {
	return nil
}

type ConstClass struct {
	NameIndex uint16
	Name *ConstUTF8
}

func NewConstClass(io io.Reader) (*ConstClass, error) {
	cc := &ConstClass{}
	if err := binary.Read(io, binary.BigEndian, &cc.NameIndex); err != nil {
		return nil, err
	}
	return cc, nil
}

func (*ConstClass) Tag() int {
	return Class
}

func (cc *ConstClass) Resolving(pool *ConstantPool) error {
	constant, err := pool.Get(cc.NameIndex)
	if err != nil {
		return err
	}
	if c, ok := constant.(*ConstUTF8); ok {
		cc.Name = c
		return nil
	}
	return ResolvingError
}

type ConstString struct {
	UTF8StringIndex uint16
	UTF8String *ConstUTF8
}

func NewConstString(io io.Reader) (*ConstString, error) {
	cs := &ConstString{}
	if err := binary.Read(io, binary.BigEndian, &cs.UTF8StringIndex); err != nil {
		return nil, err
	}
	return cs, nil
}

func (ConstString) Tag() int {
	return String
}

func (cs *ConstString) Resolving(pool *ConstantPool) error {
	constant, err := pool.Get(cs.UTF8StringIndex)
	if err != nil {
		return err
	}
	var ok bool
	if cs.UTF8String, ok = constant.(*ConstUTF8); ok {
		return nil
	}
	return ResolvingError
}

type ConstRef interface {
	GetClass() *ConstClass
	GetNameAndType() *ConstNameAndType
}

type ConstFieldRef struct {
	ClassIndex uint16
	NameAndTypeIndex uint16
	Class *ConstClass
	NameAndType *ConstNameAndType
}

func NewConstFieldRef(io io.Reader) (*ConstFieldRef, error) {
	cd := &ConstFieldRef{}
	if err := binary.Read(io, binary.BigEndian, &cd.ClassIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &cd.NameAndTypeIndex); err != nil {
		return nil, err
	}
	return cd, nil
}

func (*ConstFieldRef) Tag() int {
	return FieldRef
}

func (cd *ConstFieldRef) Resolving(pool *ConstantPool) error {
	var ok bool
	constantClass, err := pool.Get(cd.ClassIndex)
	if err != nil {
		return err
	}
	cd.Class, ok = constantClass.(*ConstClass)
	if !ok {
		return ResolvingError
	}
	constantNameAndType, err := pool.Get(cd.NameAndTypeIndex)
	if err != nil {
		return err
	}
	cd.NameAndType, ok = constantNameAndType.(*ConstNameAndType)
	if !ok {
		return ResolvingError
	}
	return nil
}

func (cd *ConstFieldRef) GetClass() *ConstClass {
	return cd.Class
}

func (cd *ConstFieldRef) GetNameAndType() *ConstNameAndType {
	return cd.NameAndType
}

type ConstMethodRef struct{*ConstFieldRef}

func (*ConstMethodRef) Tag() int {
	return MethodRef
}

func NewConstMethodRef(io io.Reader) (*ConstMethodRef, error) {
	fieldRefConst, err := NewConstFieldRef(io)
	if err != nil {
		return nil, err
	}
	return &ConstMethodRef{ConstFieldRef: fieldRefConst}, nil
}

type ConstInterfaceMethodRef struct{*ConstFieldRef}

func NewConstInterfaceMethodRef(io io.Reader) (*ConstInterfaceMethodRef, error) {
	fieldRefConst, err := NewConstFieldRef(io)
	if err != nil {
		return nil, err
	}
	return &ConstInterfaceMethodRef{fieldRefConst}, nil
}

func (*ConstInterfaceMethodRef) Tag() int {
	return InterfaceMethodRef
}

type ConstNameAndType struct {
	NameIndex uint16
	Name *ConstUTF8
	DescriptorIndex uint16
	Descriptor *ConstUTF8
}

func NewConstNameAndType(io io.Reader) (*ConstNameAndType, error) {
	cnat := &ConstNameAndType{}
	if err := binary.Read(io, binary.BigEndian, &cnat.NameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &cnat.DescriptorIndex); err != nil {
		return nil, err
	}
	return cnat, nil
}

func (*ConstNameAndType) Tag() int {
	return NameAndType
}

func (cnat *ConstNameAndType) Resolving(pool *ConstantPool) error {
	constant, err := pool.Get(cnat.NameIndex)
	if err != nil {
		return err
	}
	c, ok := constant.(*ConstUTF8)
	if !ok {
		return ResolvingError
	}
	cnat.Name = c
	constant, err = pool.Get(cnat.DescriptorIndex)
	if err != nil {
		return err
	}
	c, ok = constant.(*ConstUTF8)
	if !ok {
		return ResolvingError
	}
	cnat.Descriptor = c
	return nil
}

type ConstMethodHandle struct {
	RefKind uint8
	RefIndex uint16
	Ref ConstRef
}

func NewConstMethodHandle(io io.Reader) (*ConstMethodHandle, error) {
	cd := &ConstMethodHandle{}
	if err := binary.Read(io, binary.BigEndian, &cd.RefKind); err != nil {
		return nil, err
	}
	if cd.RefKind < RefGetField || cd.RefKind > RefInvokeInterface {
		return nil, errors.New("invaild method handle kind")
	}
	if err := binary.Read(io, binary.BigEndian, &cd.RefIndex); err != nil {
		return nil, err
	}
	return cd, nil
}

func (*ConstMethodHandle) Tag() int {
	return MethodHandle
}

func (cd *ConstMethodHandle) Resolving(pool *ConstantPool) error {
	constant, err := pool.Get(cd.RefIndex)
	if err != nil {
		return err
	}
	constRef, ok := constant.(ConstRef)
	if !ok {
		return ResolvingError
	}
	cd.Ref = constRef
	return nil
}

type ConstMethodType struct {
	DescriptorIndex uint16
	Descriptor *ConstUTF8
}

func NewConstMethodType(io io.Reader) (*ConstMethodType, error) {
	cmt := &ConstMethodType{}
	if err := binary.Read(io, binary.BigEndian, &cmt.DescriptorIndex); err != nil {
		return nil, err
	}
	return cmt, nil
}

func (ConstMethodType) Tag() int {
	return MethodType
}

func (cmt *ConstMethodType) Resolving(pool *ConstantPool) error {
	constant, err := pool.Get(cmt.DescriptorIndex)
	if err != nil {
		return err
	}
	c, ok := constant.(*ConstUTF8)
	if !ok {
		return ResolvingError
	}
	cmt.Descriptor = c
	return nil
}

type ConstInvokeDynamic struct {
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex uint16
}

func NewConstInvokeDynamic(io io.Reader) (*ConstInvokeDynamic, error) {
	cid := &ConstInvokeDynamic{}
	if err := binary.Read(io, binary.BigEndian, &cid.BootstrapMethodAttrIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &cid.NameAndTypeIndex); err != nil {
		return nil, err
	}
	return cid, nil
}

func (*ConstInvokeDynamic) Tag() int {
	return InvokeDynamic
}

func (cid *ConstInvokeDynamic) Resolving(pool *ConstantPool) error {
	return nil
}

// see java.io.DataInputStream.readUTF(DataInput)
func decodeMUTF8(bytearr []byte) (string, error) {
	utflen := len(bytearr)
	chararr := make([]uint16, utflen)

	var c, char2, char3 uint16
	count := 0
	chararr_count := 0

	for count < utflen {
		c = uint16(bytearr[count])
		if c > 127 {
			break
		}
		count++
		chararr[chararr_count] = c
		chararr_count++
	}

	for count < utflen {
		c = uint16(bytearr[count])
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			/* 0xxxxxxx*/
			count++
			chararr[chararr_count] = c
			chararr_count++
		case 12, 13:
			/* 110x xxxx   10xx xxxx*/
			count += 2
			if count > utflen {
				return "", errors.New("malformed input: partial character at end")
			}
			char2 = uint16(bytearr[count-1])
			if char2 & 0xC0 != 0x80 {
				return "", fmt.Errorf("malformed input around byte %v", count)
			}
			chararr[chararr_count] = c & 0x1F << 6 | char2 & 0x3F
			chararr_count++
		case 14:
			/* 1110 xxxx  10xx xxxx  10xx xxxx*/
			count += 3
			if count > utflen {
				return "", errors.New("malformed input: partial character at end")
			}
			char2 = uint16(bytearr[count-2])
			char3 = uint16(bytearr[count-1])
			if char2 & 0xC0 != 0x80 || char3 & 0xC0 != 0x80 {
				return "", fmt.Errorf("malformed input around byte %v", (count - 1))
			}
			chararr[chararr_count] = c & 0x0F << 12 | char2 & 0x3F << 6 | char3 & 0x3F << 0
			chararr_count++
		default:
			/* 10xx xxxx,  1111 xxxx */
			return "", fmt.Errorf("malformed input around byte %v", count)
		}
	}
	// The number of chars produced may be less than utflen
	chararr = chararr[0:chararr_count]
	runes := utf16.Decode(chararr)
	return string(runes), nil
}
