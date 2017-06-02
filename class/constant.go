package class

import (
	"io"
	"errors"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
	"math"
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
	String() string
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

func (c *ConstInteger) String() string {
	return fmt.Sprintf("ConstInteger -> %v", *c)
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

func (cf *ConstFloat) String() string {
	return fmt.Sprintf("ConstFloat -> %v", *cf)
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

func (cl *ConstLong) String() string {
	return fmt.Sprintf("ConstLong -> %v", *cl)
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

func (cl *ConstDouble) String() string {
	return fmt.Sprintf("ConstDouble -> %v", *cl)
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

func (cc *ConstClass) String() string {
	return fmt.Sprintf("ConstClass -> %v", *cc)
}

func (cc *ConstClass) Resolving(pool *ConstantPool) error {
	var err error
	if cc.Name, err = pool.GetUTF8String(cc.NameIndex); err != nil {
		return err
	}
	return nil
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
	var err error
	if cs.UTF8String, err = pool.GetUTF8String(cs.UTF8StringIndex); err != nil {
		return err
	}
	return nil
}

func (cs *ConstString) String() string {
	return fmt.Sprintf("ConstString -> %v", *cs)
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
	var err error
	if cd.Class, err = pool.GetClass(cd.ClassIndex); err != nil {
		return err
	}
	if cd.NameAndType, err = pool.GetNameAndType(cd.NameAndTypeIndex); err != nil {
		return err
	}
	return nil
}

func (cd *ConstFieldRef) GetClass() *ConstClass {
	return cd.Class
}

func (cd *ConstFieldRef) GetNameAndType() *ConstNameAndType {
	return cd.NameAndType
}

func (cd *ConstFieldRef) String() string {
	return fmt.Sprintf("ConstFieldRef -> %v", *cd)
}

type ConstMethodRef struct{*ConstFieldRef}

func (*ConstMethodRef) Tag() int {
	return MethodRef
}

func (cd *ConstMethodRef) String() string {
	return fmt.Sprintf("ConstMethodRef -> %v", *cd)
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

func (cd *ConstInterfaceMethodRef) String() string {
	return fmt.Sprintf("ConstInterfaceMethodRef -> %v", *cd)
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
	var err error
	if cnat.Name, err = pool.GetUTF8String(cnat.NameIndex); err != nil {
		return err
	}
	if cnat.Descriptor, err = pool.GetUTF8String(cnat.DescriptorIndex); err != nil {
		return err
	}
	return nil
}

func (cnat *ConstNameAndType) String() string {
	return fmt.Sprintf("ConstNameAndType -> %v", *cnat)
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

func (cd *ConstMethodHandle) String() string {
	return fmt.Sprintf("ConstMethodHandle -> %v", *cd)
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

func (cmt *ConstMethodType) String() string {
	return fmt.Sprintf("ConstMethodType -> %v", *cmt)
}

func (cmt *ConstMethodType) Resolving(pool *ConstantPool) error {
	var err error
	if cmt.Descriptor, err = pool.GetUTF8String(cmt.DescriptorIndex); err != nil {
		return err
	}
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

func (cid *ConstInvokeDynamic) String() string {
	return fmt.Sprintf("ConstInvokeDynamic -> %v", *cid)
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
