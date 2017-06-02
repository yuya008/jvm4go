package class

import (
	"io"
	"encoding/binary"
	"fmt"
)

type Annotation struct {
	Type *ConstUTF8
	ElementValPairs []*ElementValPair
}

func NewAnnotation(io io.Reader, pool *ConstantPool) (*Annotation, error) {
	var typeIndex uint16
	if err := binary.Read(io, binary.BigEndian, &typeIndex); err != nil {
		return nil, err
	}
	annotation := &Annotation{}
	var err error
	if annotation.Type, err = pool.GetUTF8String(typeIndex); err != nil {
		return nil, err
	}
	var numElementValuePairs uint16
	if err := binary.Read(io, binary.BigEndian, &numElementValuePairs); err != nil {
		return nil, err
	}
	for i := 0; i < int(numElementValuePairs); i++ {
		evp, err := NewElementValPair(io, pool)
		if err != nil {
			return nil, err
		}
		annotation.ElementValPairs = append(annotation.ElementValPairs, evp)
	}
	return annotation, nil
}

func (an *Annotation) String() string {
	return fmt.Sprintf("Annotation -> %v", *an)
}

type ElementValPair struct {
	ElementName *ConstUTF8
	Val ElementVal
}

func NewElementValPair(io io.Reader, pool *ConstantPool) (*ElementValPair, error) {
	var elementNameIndex uint16
	if err := binary.Read(io, binary.BigEndian, &elementNameIndex); err != nil {
		return nil, err
	}
	evp := &ElementValPair{}
	var err error
	if evp.ElementName, err = pool.GetUTF8String(elementNameIndex); err != nil {
		return nil, err
	}
	if evp.Val, err = NewElementVal(io, pool); err != nil {
		return nil, err
	}
	return evp, nil
}

func (evp *ElementValPair) String() string {
	return fmt.Sprintf("ElementValPair -> %v", *evp)
}

type ElementVal interface {
	Tag() uint8
	String() string
}

func NewElementVal(io io.Reader, pool *ConstantPool) (ElementVal, error) {
	var tag uint8
	if err := binary.Read(io, binary.BigEndian, &tag); err != nil {
		return nil, err
	}
	switch tag {
	case 'B':
		return NewElementValByte(io, pool)
	case 'C':
		return NewElementValChar(io, pool)
	case 'D':
		return NewElementValDouble(io, pool)
	case 'F':
		return NewElementValFloat(io, pool)
	case 'I':
		return NewElementValInt(io, pool)
	case 'J':
		return NewElementValLong(io, pool)
	case 'S':
		return NewElementValShort(io, pool)
	case 'Z':
		return NewElementValBoolean(io, pool)
	case 's':
		return NewElementValString(io, pool)
	case 'e':
		return NewElementValEnum(io, pool)
	case 'c':
		return NewElementValClass(io, pool)
	case '@':
		return NewElementValAnnotation(io, pool)
	case '[':
		return NewElementValArray(io, pool)
	}
	return nil, fmt.Errorf("unknow tag %d", tag)
}

type ElementValByte struct {
	Val *ConstInteger
}

func NewElementValByte(io io.Reader, pool *ConstantPool) (*ElementValByte, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValByte{}
	if evb.Val, err = pool.GetInteger(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValByte) Tag() uint8 {
	return 'B'
}

func (evb *ElementValByte) String() string {
	return fmt.Sprintf("ElementValByte -> %v", *evb)
}

type ElementValChar struct {
	Val *ConstInteger
}

func NewElementValChar(io io.Reader, pool *ConstantPool) (*ElementValChar, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValChar{}
	if evb.Val, err = pool.GetInteger(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValChar) Tag() uint8 {
	return 'C'
}

func (evb *ElementValChar) String() string {
	return fmt.Sprintf("ElementValChar -> %v", *evb)
}

type ElementValDouble struct {
	Val *ConstDouble
}

func NewElementValDouble(io io.Reader, pool *ConstantPool) (*ElementValDouble, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValDouble{}
	if evb.Val, err = pool.GetDouble(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValDouble) Tag() uint8 {
	return 'D'
}

func (evb *ElementValDouble) String() string {
	return fmt.Sprintf("ElementValDouble -> %v", *evb)
}

type ElementValFloat struct {
	Val *ConstFloat
}

func NewElementValFloat(io io.Reader, pool *ConstantPool) (*ElementValFloat, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValFloat{}
	if evb.Val, err = pool.GetFloat(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValFloat) Tag() uint8 {
	return 'F'
}

func (evb *ElementValFloat) String() string {
	return fmt.Sprintf("ElementValFloat -> %v", *evb)
}

type ElementValInt struct {
	Val *ConstInteger
}

func NewElementValInt(io io.Reader, pool *ConstantPool) (*ElementValInt, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValInt{}
	if evb.Val, err = pool.GetInteger(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValInt) Tag() uint8 {
	return 'I'
}

func (evb *ElementValInt) String() string {
	return fmt.Sprintf("ElementValInt -> %v", *evb)
}

type ElementValLong struct {
	Val *ConstLong
}

func NewElementValLong(io io.Reader, pool *ConstantPool) (*ElementValLong, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValLong{}
	if evb.Val, err = pool.GetLong(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValLong) Tag() uint8 {
	return 'J'
}

func (evb *ElementValLong) String() string {
	return fmt.Sprintf("ElementValLong -> %v", *evb)
}

type ElementValShort struct {
	Val *ConstInteger
}

func NewElementValShort(io io.Reader, pool *ConstantPool) (*ElementValShort, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValShort{}
	if evb.Val, err = pool.GetInteger(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValShort) Tag() uint8 {
	return 'S'
}

func (evb *ElementValShort) String() string {
	return fmt.Sprintf("ElementValShort -> %v", *evb)
}

type ElementValBoolean struct {
	Val *ConstInteger
}

func NewElementValBoolean(io io.Reader, pool *ConstantPool) (*ElementValBoolean, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValBoolean{}
	if evb.Val, err = pool.GetInteger(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValBoolean) Tag() uint8 {
	return 'Z'
}

func (evb *ElementValBoolean) String() string {
	return fmt.Sprintf("ElementValBoolean -> %v", *evb)
}

type ElementValString struct {
	Val *ConstUTF8
}

func NewElementValString(io io.Reader, pool *ConstantPool) (*ElementValString, error) {
	var valIndex uint16
	if err := binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	var err error
	evb := &ElementValString{}
	if evb.Val, err = pool.GetUTF8String(valIndex); err != nil {
		return nil, err
	}
	return evb, nil
}

func (evb *ElementValString) Tag() uint8 {
	return 's'
}

func (evb *ElementValString) String() string {
	return fmt.Sprintf("ElementValString -> %v", *evb)
}

type ElementValEnum struct {
	TypeName *ConstUTF8
	ConstName *ConstUTF8
}

func NewElementValEnum(io io.Reader, pool *ConstantPool) (*ElementValEnum, error) {
	var typeNameIndex, constNameIndex uint16
	if err := binary.Read(io, binary.BigEndian, &typeNameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &constNameIndex); err != nil {
		return nil, err
	}
	var err error
	eve := &ElementValEnum{}
	if eve.TypeName, err =  pool.GetUTF8String(typeNameIndex); err != nil {
		return nil, err
	}
	if eve.ConstName, err =  pool.GetUTF8String(constNameIndex); err != nil {
		return nil, err
	}
	return eve, nil
}

func (evb *ElementValEnum) Tag() uint8 {
	return 'e'
}

func (evb *ElementValEnum) String() string {
	return fmt.Sprintf("ElementValEnum -> %v", *evb)
}

type ElementValClass struct {
	ClassInfo *ConstUTF8
}

func NewElementValClass(io io.Reader, pool *ConstantPool) (*ElementValClass, error) {
	var classInfoIndex uint16
	if err := binary.Read(io, binary.BigEndian, &classInfoIndex); err != nil {
		return nil, err
	}
	var err error
	evc := &ElementValClass{}
	if evc.ClassInfo, err = pool.GetUTF8String(classInfoIndex); err != nil {
		return nil, err
	}
	return evc, nil
}

func (evb *ElementValClass) Tag() uint8 {
	return 'c'
}

func (evb *ElementValClass) String() string {
	return fmt.Sprintf("ElementValClass -> %v", *evb)
}

type ElementValAnnotation struct {
	AnnotationVal *Annotation
}

func NewElementValAnnotation(io io.Reader, pool *ConstantPool) (*ElementValAnnotation, error) {
	eva := &ElementValAnnotation{}
	var err error
	if eva.AnnotationVal, err = NewAnnotation(io, pool); err != nil {
		return nil, err
	}
	return eva, nil
}

func (evb *ElementValAnnotation) Tag() uint8 {
	return '@'
}

func (evb *ElementValAnnotation) String() string {
	return fmt.Sprintf("ElementValAnnotation -> %v", *evb)
}

type ElementValArray struct {
	Val []ElementVal
}

func NewElementValArray(io io.Reader, pool *ConstantPool) (*ElementValArray, error) {
	var numValues uint16
	if err := binary.Read(io, binary.BigEndian, &numValues); err != nil {
		return nil, err
	}
	eva := &ElementValArray{}
	for i := 0; i < int(numValues); i++ {
		ev, err := NewElementVal(io, pool)
		if err != nil {
			return nil, err
		}
		eva.Val = append(eva.Val, ev)
	}
	return eva, nil
}

func (evb *ElementValArray) Tag() uint8 {
	return '['
}

func (evb *ElementValArray) String() string {
	return fmt.Sprintf("ElementValArray -> %v", *evb)
}

type TypeAnnotation struct {
	TargetType uint8
	Target Target
	TargetPath *TypePath
	Type *ConstUTF8
	ElementValPairs []*ElementValPair
}

func NewTypeAnnotation(io io.Reader, pool *ConstantPool) (*TypeAnnotation, error) {
	ta := &TypeAnnotation{}
	if err := binary.Read(io, binary.BigEndian, &ta.TargetType); err != nil {
		return nil, err
	}
	var err error
	if ta.Target, err = NewTarget(ta.TargetType, io); err != nil {
		return nil, err
	}
	if ta.TargetPath, err = NewTypePath(io); err != nil {
		return nil, err
	}
	annotation, err := NewAnnotation(io, pool)
	if err != nil {
		return nil, err
	}
	ta.Type = annotation.Type
	ta.ElementValPairs = annotation.ElementValPairs
	return ta, nil
}

func (t *TypeAnnotation) String() string {
	return fmt.Sprintf("TypeAnnotation -> %v", *t)
}

type Target interface {}
func NewTarget(targetType uint8, io io.Reader) (Target, error) {
	switch targetType {
	case 0x00, 0x01: // type_parameter_target
		return NewTypeParameterTarget(io)
	case 0x10: // supertype_target
		return NewSupertypeTarget(io)
	case 0x11, 0x12: // type_parameter_bound_target
		return NewTypeParameterBoundTarget(io)
	case 0x13, 0x14, 0x15: // empty_target
		return NewEmptyTarget()
	case 0x16: // formal_parameter_target
		return NewFormalParameterTarget(io)
	case 0x17: // throws_target
		return NewThrowsTarget(io)
	case 0x40, 0x41: // localvar_target
		return NewLocalvarTarget(io)
	case 0x42: // catch_target
		return NewCatchTarget(io)
	case 0x43, 0x44, 0x45, 0x46: // offset_target
		return NewOffsetTarget(io)
	case 0x47, 0x48, 0x49, 0x4A, 0x4B: // type_argument_target
		return NewTypeArgumentTarget(io)
	}
	return nil, fmt.Errorf("unknow target %d", targetType)
}

type TypeParameterTarget struct {
	TypeParameterIndex uint8
}

func NewTypeParameterTarget(io io.Reader) (*TypeParameterTarget, error) {
	tpt := &TypeParameterTarget{}
	if err := binary.Read(io, binary.BigEndian, &tpt.TypeParameterIndex); err != nil {
		return nil, err
	}
	return tpt, nil
}

type SupertypeTarget struct {
	SupertypeIndex uint16
}

func NewSupertypeTarget(io io.Reader) (*SupertypeTarget, error) {
	st := &SupertypeTarget{}
	if err := binary.Read(io, binary.BigEndian, &st.SupertypeIndex); err != nil {
		return nil, err
	}
	return st, nil
}

type TypeParameterBoundTarget struct {
	TypeParameterIndex, BoundIndex uint8
}

func NewTypeParameterBoundTarget(io io.Reader) (*TypeParameterBoundTarget, error) {
	tpbt := &TypeParameterBoundTarget{}
	if err := binary.Read(io, binary.BigEndian, &tpbt.TypeParameterIndex); err != nil {
		return nil, err
	}
	return tpbt, nil
}

type EmptyTarget struct {}

func NewEmptyTarget() (*EmptyTarget, error) {
	return &EmptyTarget{}, nil
}

type FormalParameterTarget struct {
	FormalParameterIndex uint8
}

func NewFormalParameterTarget(io io.Reader) (*FormalParameterTarget, error) {
	fpt := &FormalParameterTarget{}
	if err := binary.Read(io, binary.BigEndian, &fpt.FormalParameterIndex); err != nil {
		return nil, err
	}
	return fpt, nil
}

type ThrowsTarget struct {
	ThrowsTypeIndex uint8
}

func NewThrowsTarget(io io.Reader) (*ThrowsTarget, error) {
	tt := &ThrowsTarget{}
	if err := binary.Read(io, binary.BigEndian, &tt.ThrowsTypeIndex); err != nil {
		return nil, err
	}
	return tt, nil
}

type Table struct {
	StartPC, Length, Index uint16
}

func NewTable(io io.Reader) (*Table, error) {
	t := &Table{}
	if err := binary.Read(io, binary.BigEndian, &t.StartPC); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &t.Length); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &t.Index); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Table) String() string {
	return fmt.Sprintf("Table -> %v", *t)
}

type LocalvarTarget struct {
	Table []*Table
}

func NewLocalvarTarget(io io.Reader) (*LocalvarTarget, error) {
	var tableLen uint16
	if err := binary.Read(io, binary.BigEndian, &tableLen); err != nil {
		return nil, err
	}
	lt := &LocalvarTarget{}
	for i := 0; i < int(tableLen); i++ {
		table, err := NewTable(io)
		if err != nil {
			return nil, err
		}
		lt.Table = append(lt.Table, table)
	}
	return lt, nil
}

type CatchTarget struct {
	ExceptionTableIndex uint16
}

func NewCatchTarget(io io.Reader) (*CatchTarget, error) {
	ct := &CatchTarget{}
	if err := binary.Read(io, binary.BigEndian, &ct.ExceptionTableIndex); err != nil {
		return nil, err
	}
	return ct, nil
}

type OffsetTarget struct {
	Offset uint16
}

func NewOffsetTarget(io io.Reader) (*OffsetTarget, error) {
	ot := &OffsetTarget{}
	if err := binary.Read(io, binary.BigEndian, &ot.Offset); err != nil {
		return nil, err
	}
	return ot, nil
}

type TypeArgumentTarget struct {
	Offset uint16
	TypeArgumentIndex uint8
}

func NewTypeArgumentTarget(io io.Reader) (*TypeArgumentTarget, error) {
	tat := &TypeArgumentTarget{}
	if err := binary.Read(io, binary.BigEndian, &tat.Offset); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &tat.TypeArgumentIndex); err != nil {
		return nil, err
	}
	return tat, nil
}

type Path struct {
	TypePathKind, TypeArgumentIndex uint8
}

func NewPath(io io.Reader) (*Path, error) {
	p := &Path{}
	if err := binary.Read(io, binary.BigEndian, &p.TypePathKind); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &p.TypeArgumentIndex); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Path) String() string {
	return fmt.Sprintf("Path -> %v", *p)
}

type TypePath struct {
	Path []*Path
}

func NewTypePath(io io.Reader) (*TypePath, error) {
	var pathLength uint8
	if err := binary.Read(io, binary.BigEndian, &pathLength); err != nil {
		return nil, err
	}
	tp := &TypePath{}
	for i := 0; i < int(pathLength); i++ {
		path, err := NewPath(io)
		if err != nil {
			return nil, err
		}
		tp.Path = append(tp.Path, path)
	}
	return tp, nil
}