package class

import (
	"io"
	"encoding/binary"
	"fmt"
)

const (
	ConstantValue = "ConstantValue"
	Code = "Code"
	Exceptions = "Exceptions"
	SourceFile = "SourceFile"
	LineNumberTable = "LineNumberTable"
	LocalVariableTable = "LocalVariableTable"
	InnerClasses = "InnerClasses"
	Synthetic = "Synthetic"
	Deprecated = "Deprecated"
	EnclosingMethod = "EnclosingMethod"
	Signature = "Signature"
	SourceDebugExtension = "SourceDebugExtension"
	LocalVariableTypeTable = "LocalVariableTypeTable"
	RuntimeVisibleAnnotations = "RuntimeVisibleAnnotations"
	RuntimeInvisibleAnnotations = "RuntimeInvisibleAnnotations"
	RuntimeVisibleParameterAnnotations = "RuntimeVisibleParameterAnnotations"
	RuntimeInvisibleParameterAnnotations = "RuntimeInvisibleParameterAnnotations"
	AnnotationDefault = "AnnotationDefault"
	StackMapTable = "StackMapTable"
	BootstrapMethods = "BootstrapMethods"
	RuntimeVisibleTypeAnnotations = "RuntimeVisibleTypeAnnotations"
	RuntimeInvisibleTypeAnnotations = "RuntimeInvisibleTypeAnnotations"
	MethodParameters = "MethodParameters"
)

type Attr interface {
	Name() string
	String() string
}

func ReadAttr(reader io.Reader, pool *ConstantPool) (Attr, error) {
	var attrNameIndex uint16
	if err := binary.Read(reader, binary.BigEndian, &attrNameIndex); err != nil {
		return nil, err
	}
	attrName, err := pool.GetUTF8String(attrNameIndex)
	if err != nil {
		return nil, err
	}
	var attrLength uint32
	if err := binary.Read(reader, binary.BigEndian, &attrLength); err != nil {
		return nil, err
	}
	switch attrName.String() {
	case ConstantValue:
		return NewAttrConstantValue(reader, pool)
	case Code:
		return NewAttrCode(reader, pool)
	case StackMapTable:
		return NewAttrStackMapTable(reader, pool)
	case AnnotationDefault:
		return NewAttrAnnotationDefault(reader, pool)
	case BootstrapMethods:
		return NewAttrBootstrapMethods(reader, pool)
	case Deprecated:
		return NewAttrDeprecated()
	case EnclosingMethod:
		return NewAttrEnclosingMethod(reader, pool)
	case Exceptions:
		return NewAttrExceptions(reader, pool)
	case InnerClasses:
		return NewAttrInnerClasses(reader, pool)
	case LineNumberTable:
		return NewAttrLineNumberTable(reader)
	case LocalVariableTable:
		return NewAttrLocalVariableTable(reader, pool)
	case LocalVariableTypeTable:
		return NewAttrLocalVariableTypeTable(reader, pool)
	case MethodParameters:
		return NewAttrMethodParameters(reader, pool)
	case RuntimeInvisibleAnnotations:
		return NewAttrRuntimeInvisibleAnnotations(reader, pool)
	case RuntimeInvisibleParameterAnnotations:
		return NewAttrRuntimeInvisibleParameterAnnotations(reader, pool)
	case RuntimeInvisibleTypeAnnotations:
		return NewAttrRuntimeInvisibleTypeAnnotations(reader, pool)
	case RuntimeVisibleAnnotations:
		return NewAttrRuntimeVisibleAnnotations(reader, pool)
	case RuntimeVisibleParameterAnnotations:
		return NewAttrRuntimeVisibleParameterAnnotations(reader, pool)
	case RuntimeVisibleTypeAnnotations:
		return NewAttrRuntimeVisibleTypeAnnotations(reader, pool)
	case Signature:
		return NewAttrSignature(reader, pool)
	case SourceDebugExtension:
		return NewAttrSourceDebugExtension(attrLength, reader)
	case SourceFile:
		return NewAttrSourceFile(reader, pool)
	case Synthetic:
		return NewAttrSynthetic()
	}
	return nil, fmt.Errorf("%s unknow name", attrName.String())
}

type AttrConstantValue struct {
	Val Constant
}

func NewAttrConstantValue(io io.Reader, pool *ConstantPool) (*AttrConstantValue, error) {
	var valIndex uint16
	var err error
	if err = binary.Read(io, binary.BigEndian, &valIndex); err != nil {
		return nil, err
	}
	attr := &AttrConstantValue{}
	if attr.Val, err = pool.Get(valIndex); err != nil {
		return nil, err
	}
	return attr, nil
}

func (attr *AttrConstantValue) Name() string {
	return ConstantValue
}

func (attr *AttrConstantValue) String() string {
	return attr.Val.String()
}

type Exception struct {
	StartPC uint16
	EndPC uint16
	HandlerPC uint16
	CatchType *ConstClass
}

func NewException(io io.Reader, pool *ConstantPool) (*Exception, error) {
	ex := &Exception{}
	if err := binary.Read(io, binary.BigEndian, &ex.StartPC); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &ex.EndPC); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &ex.HandlerPC); err != nil {
		return nil, err
	}
	var catchTypeIndex uint16
	if err := binary.Read(io, binary.BigEndian, &catchTypeIndex); err != nil {
		return nil, err
	}
	var err error
	if ex.CatchType, err = pool.GetClass(catchTypeIndex); err != nil {
		return nil, err
	}
	return ex, nil
}

type AttrCode struct {
	MaxStack uint16
	MaxLocals uint16
	Code []byte
	ExceptionTable []*Exception
	Attrs []Attr
}

func NewAttrCode(io io.Reader, pool *ConstantPool) (*AttrCode, error) {
	code := &AttrCode{}
	if err := binary.Read(io, binary.BigEndian, &code.MaxStack); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &code.MaxLocals); err != nil {
		return nil, err
	}
	var codeLength uint32
	if err := binary.Read(io, binary.BigEndian, &codeLength); err != nil {
		return nil, err
	}
	code.Code = make([]byte, codeLength)
	if err := binary.Read(io, binary.BigEndian, &code.Code); err != nil {
		return nil, err
	}
	var exceptionLen uint16
	if err := binary.Read(io, binary.BigEndian, &exceptionLen); err != nil {
		return nil, err
	}
	for i := 0; i < int(exceptionLen); i++ {
		ex, err := NewException(io, pool)
		if err != nil {
			return nil, err
		}
		code.ExceptionTable = append(code.ExceptionTable, ex)
	}
	var attrCount uint16
	if err := binary.Read(io, binary.BigEndian, &attrCount); err != nil {
		return nil, err
	}
	for i := 0; i < int(attrCount); i++ {
		attr, err := ReadAttr(io, pool)
		if err != nil {
			return nil, err
		}
		code.Attrs = append(code.Attrs, attr)
	}
	return code, nil
}

func (c *AttrCode) Name() string {
	return Code
}

func (attr *AttrCode) String() string {
	return fmt.Sprintf("AttrCode ->  %v", *attr)
}

type AttrStackMapTable struct {
	Entries []StackMapFrame
}

func NewAttrStackMapTable(io io.Reader, pool *ConstantPool) (*AttrStackMapTable, error) {
	stackMapTable := &AttrStackMapTable{}
	var numberOfEntries uint16
	if err := binary.Read(io, binary.BigEndian, &numberOfEntries); err != nil {
		return nil, err
	}
	for i := 0; i < int(numberOfEntries); i++ {
		s, err := NewStackMapFrame(io, pool)
		if err != nil {
			return nil, err
		}
		stackMapTable.Entries = append(stackMapTable.Entries, s)
	}
	return stackMapTable, nil
}

func (c *AttrStackMapTable) Name() string {
	return StackMapTable
}

func (c *AttrStackMapTable) String() string {
	return fmt.Sprintf("AttrStackMapTable ->  %v", *c)
}

type AttrExceptions struct {
	ExceptionTable []*ConstClass
}

func NewAttrExceptions(io io.Reader, pool *ConstantPool) (*AttrExceptions, error) {
	var numberOfExceptions uint16
	if err := binary.Read(io, binary.BigEndian, &numberOfExceptions); err != nil {
		return nil, err
	}
	ex := &AttrExceptions{}
	var exIndex uint16
	for i := 0; i < int(numberOfExceptions); i++ {
		if err := binary.Read(io, binary.BigEndian, &exIndex); err != nil {
			return nil, err
		}
		exEntry, err := pool.GetClass(exIndex)
		if err != nil {
			return nil, err
		}
		ex.ExceptionTable = append(ex.ExceptionTable, exEntry)
	}
	return ex, nil
}

func (c *AttrExceptions) Name() string {
	return Exceptions
}

func (c *AttrExceptions) String() string {
	return fmt.Sprintf("AttrExceptions ->  %v", *c)
}

type Classes struct {
	InnerClass *ConstClass
	OuterClass *ConstClass
	InnerName *ConstUTF8
	InnerClassAccessFlags uint16
}

func NewClasses(io io.Reader, pool *ConstantPool) (*Classes, error) {
	classes := &Classes{}
	var innerClassIndex, outerClassIndex, innerNameIndex uint16
	if err := binary.Read(io, binary.BigEndian, &innerClassIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &outerClassIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &innerNameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &classes.InnerClassAccessFlags); err != nil {
		return nil, err
	}
	var err error
	if classes.InnerClass, err = pool.GetClass(innerClassIndex); err != nil {
		return nil, err
	}
	if classes.OuterClass, err = pool.GetClass(outerClassIndex); err != nil {
		return nil, err
	}
	if classes.InnerName, err = pool.GetUTF8String(innerNameIndex); err != nil {
		return nil, err
	}
	return classes, nil
}

type AttrInnerClasses struct {
	Classes []*Classes
}

func NewAttrInnerClasses(io io.Reader, pool *ConstantPool) (*AttrInnerClasses, error) {
	var numberOfClasses uint16
	if err := binary.Read(io, binary.BigEndian, &numberOfClasses); err != nil {
		return nil, err
	}
	a := &AttrInnerClasses{}
	for i := 0; i < int(numberOfClasses); i++ {
		classes, err := NewClasses(io, pool)
		if err != nil {
			return nil, err
		}
		a.Classes = append(a.Classes, classes)
	}
	return a, nil
}

func (c *AttrInnerClasses) Name() string {
	return InnerClasses
}

func (c *AttrInnerClasses) String() string {
	return fmt.Sprintf("AttrInnerClasses ->  %v", *c)
}

type AttrEnclosingMethod struct {
	Class *ConstClass
	Method *ConstNameAndType
}

func NewAttrEnclosingMethod(io io.Reader, pool *ConstantPool) (*AttrEnclosingMethod, error) {
	var classIndex, methodIndex uint16
	var err error
	if err := binary.Read(io, binary.BigEndian, &classIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &methodIndex); err != nil {
		return nil, err
	}
	a := &AttrEnclosingMethod{}
	if a.Class, err = pool.GetClass(classIndex); err != nil {
		return nil, err
	}
	if a.Method, err = pool.GetNameAndType(methodIndex); err != nil {
		return nil, err
	}
	return a, nil
}

func (c *AttrEnclosingMethod) Name() string {
	return EnclosingMethod
}

func (c *AttrEnclosingMethod) String() string {
	return fmt.Sprintf("AttrEnclosingMethod ->  %v", *c)
}

type AttrSynthetic struct {}

func NewAttrSynthetic() (*AttrSynthetic, error) {
	return &AttrSynthetic{}, nil
}

func (c *AttrSynthetic) Name() string {
	return Synthetic
}

func (c *AttrSynthetic) String() string {
	return fmt.Sprint("AttrSynthetic -> {}")
}

type AttrSignature struct {
	Signature *ConstUTF8
}

func NewAttrSignature(io io.Reader, pool *ConstantPool) (*AttrSignature, error) {
	s := &AttrSignature{}
	var signatureIndex uint16
	var err error
	if err := binary.Read(io, binary.BigEndian, &signatureIndex); err != nil {
		return nil, err
	}
	if s.Signature, err = pool.GetUTF8String(signatureIndex); err != nil {
		return nil, err
	}
	return s, nil
}

func (c *AttrSignature) Name() string {
	return Signature
}

func (c *AttrSignature) String() string {
	return fmt.Sprintf("AttrSignature ->  %v", *c)
}

type AttrSourceFile struct {
	SourceFile *ConstUTF8
}

func NewAttrSourceFile(io io.Reader, pool *ConstantPool) (*AttrSourceFile, error) {
	var sourceFileIndex uint16
	var err error
	if err := binary.Read(io, binary.BigEndian, &sourceFileIndex); err != nil {
		return nil, err
	}
	sf := &AttrSourceFile{}
	if sf.SourceFile, err = pool.GetUTF8String(sourceFileIndex); err != nil {
		return nil, err
	}
	return sf, nil
}

func (c *AttrSourceFile) Name() string {
	return SourceFile
}

func (c *AttrSourceFile) String() string {
	return fmt.Sprintf("AttrSourceFile ->  %v", *c)
}

type AttrSourceDebugExtension struct {
	DebugExtension []byte
}

func NewAttrSourceDebugExtension(len uint32, io io.Reader) (*AttrSourceDebugExtension, error) {
	sde := &AttrSourceDebugExtension{DebugExtension:make([]byte, len)}
	if err := binary.Read(io, binary.BigEndian, sde.DebugExtension); err != nil {
		return nil, err
	}
	return sde, nil
}

func (c *AttrSourceDebugExtension) Name() string {
	return SourceDebugExtension
}

func (c *AttrSourceDebugExtension) String() string {
	return fmt.Sprintf("AttrSourceDebugExtension -> {%s}", string(c.DebugExtension))
}

type LineNumberTableEntry struct {
	StartPC, LineNumber uint16
}

func NewLineNumberTableEntry(io io.Reader) (*LineNumberTableEntry, error) {
	l := &LineNumberTableEntry{}
	if err := binary.Read(io, binary.BigEndian, &l.StartPC); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &l.LineNumber); err != nil {
		return nil, err
	}
	return l, nil
}

func (lnte *LineNumberTableEntry) String() string {
	return fmt.Sprintf("LineNumberTableEntry -> %v", *lnte)
}

type AttrLineNumberTable struct {
	LineNumberTable []*LineNumberTableEntry
}

func NewAttrLineNumberTable(io io.Reader) (*AttrLineNumberTable, error) {
	var lineNumberTableLen uint16
	if err := binary.Read(io, binary.BigEndian, &lineNumberTableLen); err != nil {
		return nil, err
	}
	lnt := &AttrLineNumberTable{}
	for i := 0; i < int(lineNumberTableLen); i++ {
		entry, err := NewLineNumberTableEntry(io)
		if err != nil {
			return nil, err
		}
		lnt.LineNumberTable = append(lnt.LineNumberTable, entry)
	}
	return lnt, nil
}

func (c *AttrLineNumberTable) Name() string {
	return LineNumberTable
}

func (c *AttrLineNumberTable) String() string {
	return fmt.Sprintf("AttrLineNumberTable ->  %v", *c)
}

type LocalVarTableEntry struct {
	StartPC uint16
	Length uint16
	Name *ConstUTF8
	Descriptor *ConstUTF8
	Index uint16
}

func NewLocalVarTableEntry(io io.Reader, pool *ConstantPool) (*LocalVarTableEntry, error) {
	vte := &LocalVarTableEntry{}
	if err := binary.Read(io, binary.BigEndian, &vte.StartPC); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &vte.Length); err != nil {
		return nil, err
	}
	var nameIndex, descriptorIndex uint16
	if err := binary.Read(io, binary.BigEndian, &nameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &descriptorIndex); err != nil {
		return nil, err
	}
	var err error
	if vte.Name, err = pool.GetUTF8String(nameIndex); err != nil {
		return nil, err
	}
	if vte.Descriptor, err = pool.GetUTF8String(descriptorIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &vte.Index); err != nil {
		return nil, err
	}
	return vte, nil
}

type AttrLocalVariableTable struct {
	LocalVarTable []*LocalVarTableEntry
}

func NewAttrLocalVariableTable(io io.Reader, pool *ConstantPool) (*AttrLocalVariableTable, error) {
	var localVarTableEntryLen uint16
	if err := binary.Read(io, binary.BigEndian, &localVarTableEntryLen); err != nil {
		return nil, err
	}
	lvt := &AttrLocalVariableTable{}
	for i := 0; i < int(localVarTableEntryLen); i++ {
		entry, err := NewLocalVarTableEntry(io, pool)
		if err != nil {
			return nil, err
		}
		lvt.LocalVarTable = append(lvt.LocalVarTable, entry)
	}
	return lvt, nil
}

func (c *AttrLocalVariableTable) Name() string {
	return LocalVariableTable
}

func (c *AttrLocalVariableTable) String() string {
	return fmt.Sprintf("AttrLocalVariableTable ->  %v", *c)
}

type LocalVarTypeTableEntry struct {
	StartPC uint16
	Length uint16
	Name *ConstUTF8
	Signature *ConstUTF8
	Index uint16
}

func NewLocalVarTypeTableEntry(io io.Reader, pool *ConstantPool) (*LocalVarTypeTableEntry, error) {
	lvtt := &LocalVarTypeTableEntry{}
	if err := binary.Read(io, binary.BigEndian, &lvtt.StartPC); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &lvtt.Length); err != nil {
		return nil, err
	}
	var nameIndex, signatureIndex uint16
	if err := binary.Read(io, binary.BigEndian, &nameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &signatureIndex); err != nil {
		return nil, err
	}
	var err error
	if lvtt.Name, err = pool.GetUTF8String(nameIndex); err != nil {
		return nil, err
	}
	if lvtt.Signature, err = pool.GetUTF8String(signatureIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &lvtt.Index); err != nil {
		return nil, err
	}
	return lvtt, nil
}

type AttrLocalVariableTypeTable struct {
	LocalVarTypeTable []*LocalVarTypeTableEntry
}

func NewAttrLocalVariableTypeTable(io io.Reader, pool *ConstantPool) (*AttrLocalVariableTypeTable, error) {
	var len uint16
	if err := binary.Read(io, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	lvtt := &AttrLocalVariableTypeTable{}
	for i := 0; i < int(len); i++ {
		lvtte, err := NewLocalVarTypeTableEntry(io, pool)
		if err != nil {
			return nil, err
		}
		lvtt.LocalVarTypeTable = append(lvtt.LocalVarTypeTable, lvtte)
	}
	return lvtt, nil
}

func (c *AttrLocalVariableTypeTable) Name() string {
	return LocalVariableTypeTable
}

func (c *AttrLocalVariableTypeTable) String() string {
	return fmt.Sprintf("AttrLocalVariableTypeTable ->  %v", *c)
}

type AttrDeprecated struct {}

func NewAttrDeprecated() (*AttrDeprecated, error) {
	return &AttrDeprecated{}, nil
}

func (c *AttrDeprecated) Name() string {
	return Deprecated
}

func (c *AttrDeprecated) String() string {
	return fmt.Sprint("AttrDeprecated -> {}")
}

type AttrRuntimeVisibleAnnotations struct {
	Annotations []*Annotation
}

func NewAttrRuntimeVisibleAnnotations(io io.Reader, pool *ConstantPool) (*AttrRuntimeVisibleAnnotations, error) {
	var n uint16
	if err := binary.Read(io, binary.BigEndian, &n); err != nil {
		return nil, err
	}
	rva := &AttrRuntimeVisibleAnnotations{}
	for i := 0; i < int(n); i++ {
		annotation, err := NewAnnotation(io, pool)
		if err != nil {
			return nil, err
		}
		rva.Annotations = append(rva.Annotations, annotation)
	}
	return rva, nil
}

func (c *AttrRuntimeVisibleAnnotations) Name() string {
	return RuntimeVisibleAnnotations
}

func (c *AttrRuntimeVisibleAnnotations) String() string {
	return fmt.Sprintf("AttrRuntimeVisibleAnnotations ->  %v", *c)
}

type AttrRuntimeInvisibleAnnotations struct {
	Annotations []*Annotation
}

func NewAttrRuntimeInvisibleAnnotations(io io.Reader, pool *ConstantPool) (*AttrRuntimeInvisibleAnnotations, error) {
	var n uint16
	if err := binary.Read(io, binary.BigEndian, &n); err != nil {
		return nil, err
	}
	rva := &AttrRuntimeInvisibleAnnotations{}
	for i := 0; i < int(n); i++ {
		annotation, err := NewAnnotation(io, pool)
		if err != nil {
			return nil, err
		}
		rva.Annotations = append(rva.Annotations, annotation)
	}
	return rva, nil
}

func (c *AttrRuntimeInvisibleAnnotations) Name() string {
	return RuntimeInvisibleAnnotations
}

func (c *AttrRuntimeInvisibleAnnotations) String() string {
	return fmt.Sprintf("AttrRuntimeInvisibleAnnotations ->  %v", *c)
}

type ParameterAnnotation struct {
	Annotations []*Annotation
}

func NewParameterAnnotation(io io.Reader, pool *ConstantPool) (*ParameterAnnotation, error) {
	var numAnnotations uint16
	if err := binary.Read(io, binary.BigEndian, &numAnnotations); err != nil {
		return nil, err
	}
	pa := &ParameterAnnotation{}
	for i := 0; i < int(numAnnotations); i++ {
		annotation, err := NewAnnotation(io, pool)
		if err != nil {
			return nil, err
		}
		pa.Annotations = append(pa.Annotations, annotation)
	}
	return pa, nil
}

func (pa *ParameterAnnotation) String() string {
	return fmt.Sprintf("ParameterAnnotation -> %v", *pa)
}

type AttrRuntimeVisibleParameterAnnotations struct {
	ParameterAnnotations []*ParameterAnnotation
}

func NewAttrRuntimeVisibleParameterAnnotations(io io.Reader, pool *ConstantPool) (*AttrRuntimeVisibleParameterAnnotations, error) {
	var numParameters uint8
	if err := binary.Read(io, binary.BigEndian, &numParameters); err != nil {
		return nil, err
	}
	rvpa := &AttrRuntimeVisibleParameterAnnotations{}
	for i := 0; i < int(numParameters); i++ {
		pa, err := NewParameterAnnotation(io, pool)
		if err != nil {
			return nil, err
		}
		rvpa.ParameterAnnotations = append(rvpa.ParameterAnnotations, pa)
	}
	return rvpa, nil
}

func (c *AttrRuntimeVisibleParameterAnnotations) Name() string {
	return RuntimeVisibleParameterAnnotations
}

func (c *AttrRuntimeVisibleParameterAnnotations) String() string {
	return fmt.Sprintf("AttrRuntimeVisibleParameterAnnotations ->  %v", *c)
}

type AttrRuntimeInvisibleParameterAnnotations struct {
	ParameterAnnotations []*ParameterAnnotation
}

func NewAttrRuntimeInvisibleParameterAnnotations(io io.Reader, pool *ConstantPool) (*AttrRuntimeInvisibleParameterAnnotations, error) {
	var numParameters uint8
	if err := binary.Read(io, binary.BigEndian, &numParameters); err != nil {
		return nil, err
	}
	rvpa := &AttrRuntimeInvisibleParameterAnnotations{}
	for i := 0; i < int(numParameters); i++ {
		pa, err := NewParameterAnnotation(io, pool)
		if err != nil {
			return nil, err
		}
		rvpa.ParameterAnnotations = append(rvpa.ParameterAnnotations, pa)
	}
	return rvpa, nil
}

func (c *AttrRuntimeInvisibleParameterAnnotations) Name() string {
	return RuntimeInvisibleParameterAnnotations
}

func (c *AttrRuntimeInvisibleParameterAnnotations) String() string {
	return fmt.Sprintf("AttrRuntimeInvisibleParameterAnnotations ->  %v", *c)
}

type AttrRuntimeVisibleTypeAnnotations struct {
	Annotations []*TypeAnnotation
}

func NewAttrRuntimeVisibleTypeAnnotations(io io.Reader, pool *ConstantPool) (*AttrRuntimeVisibleTypeAnnotations, error) {
	var numAnnotations uint16
	if err := binary.Read(io, binary.BigEndian, &numAnnotations); err != nil {
		return nil, err
	}
	rvta := &AttrRuntimeVisibleTypeAnnotations{}
	for i := 0; i < int(numAnnotations); i++ {
		typeAnnotation, err := NewTypeAnnotation(io, pool)
		if err != nil {
			return nil, err
		}
		rvta.Annotations = append(rvta.Annotations, typeAnnotation)
	}
	return rvta, nil
}

func (c *AttrRuntimeVisibleTypeAnnotations) Name() string {
	return RuntimeVisibleTypeAnnotations
}

func (c *AttrRuntimeVisibleTypeAnnotations) String() string {
	return fmt.Sprintf("AttrRuntimeVisibleTypeAnnotations ->  %v", *c)
}

type AttrRuntimeInvisibleTypeAnnotations struct {
	Annotations []*TypeAnnotation
}

func NewAttrRuntimeInvisibleTypeAnnotations(io io.Reader, pool *ConstantPool) (*AttrRuntimeInvisibleTypeAnnotations, error) {
	var numAnnotations uint16
	if err := binary.Read(io, binary.BigEndian, &numAnnotations); err != nil {
		return nil, err
	}
	rvta := &AttrRuntimeInvisibleTypeAnnotations{}
	for i := 0; i < int(numAnnotations); i++ {
		typeAnnotation, err := NewTypeAnnotation(io, pool)
		if err != nil {
			return nil, err
		}
		rvta.Annotations = append(rvta.Annotations, typeAnnotation)
	}
	return rvta, nil
}

func (c *AttrRuntimeInvisibleTypeAnnotations) Name() string {
	return RuntimeInvisibleTypeAnnotations
}

func (c *AttrRuntimeInvisibleTypeAnnotations) String() string {
	return fmt.Sprintf("AttrRuntimeInvisibleTypeAnnotations ->  %v", *c)
}

type AttrAnnotationDefault struct {
	DefaultVal ElementVal
}

func NewAttrAnnotationDefault(io io.Reader, pool *ConstantPool) (*AttrAnnotationDefault, error) {
	ad := &AttrAnnotationDefault{}
	var err error
	if ad.DefaultVal, err = NewElementVal(io, pool); err != nil {
		return nil, err
	}
	return ad, nil
}

func (c *AttrAnnotationDefault) Name() string {
	return AnnotationDefault
}

func (c *AttrAnnotationDefault) String() string {
	return fmt.Sprintf("AttrAnnotationDefault ->  %v", *c)
}

type BootstrapMethod struct {
	BootstrapMethodRef *ConstMethodHandle
	BootstrapArguments []Constant
}

func NewBootstrapMethod(io io.Reader, pool *ConstantPool) (*BootstrapMethod, error) {
	var err error
	var bootstrapMethodRef uint16
	if err := binary.Read(io, binary.BigEndian, &bootstrapMethodRef); err != nil {
		return nil, err
	}
	bsm := &BootstrapMethod{}
	if bsm.BootstrapMethodRef, err = pool.GetMethodHandle(bootstrapMethodRef); err != nil {
		return nil, err
	}
	var numBootstrapArguments, bootstrapArgumentsIndex uint16
	if err := binary.Read(io, binary.BigEndian, &numBootstrapArguments); err != nil {
		return nil, err
	}
	for i := 0; i < int(numBootstrapArguments); i++ {
		if err := binary.Read(io, binary.BigEndian, &bootstrapArgumentsIndex); err != nil {
			return nil, err
		}
		constant, err := pool.Get(bootstrapArgumentsIndex)
		if err != nil {
			return nil, err
		}
		bsm.BootstrapArguments = append(bsm.BootstrapArguments, constant)
	}
	return bsm, nil
}

type AttrBootstrapMethods struct {
	BootstrapMethods []*BootstrapMethod
}

func NewAttrBootstrapMethods(io io.Reader, pool *ConstantPool) (*AttrBootstrapMethods, error) {
	var numBootstrapMethods uint16
	if err := binary.Read(io, binary.BigEndian, &numBootstrapMethods); err != nil {
		return nil, err
	}
	bm := &AttrBootstrapMethods{}
	for i := 0; i < int(numBootstrapMethods); i++ {
		bsm, err := NewBootstrapMethod(io, pool)
		if err != nil {
			return nil, err
		}
		bm.BootstrapMethods = append(bm.BootstrapMethods, bsm)
	}
	return bm, nil
}

func (c *AttrBootstrapMethods) Name() string {
	return BootstrapMethods
}

func (c *AttrBootstrapMethods) String() string {
	return fmt.Sprintf("AttrBootstrapMethods ->  %v", *c)
}

type Parameter struct {
	Name *ConstUTF8
	AccessFlags uint16
}

func NewParameter(io io.Reader, pool *ConstantPool) (*Parameter, error) {
	var nameIndex uint16
	if err := binary.Read(io, binary.BigEndian, &nameIndex); err != nil {
		return nil, err
	}
	var err error
	p := &Parameter{}
	if p.Name, err = pool.GetUTF8String(nameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(io, binary.BigEndian, &p.AccessFlags); err != nil {
		return nil, err
	}
	return p, nil
}

type AttrMethodParameters struct {
	Parameters []*Parameter
}

func NewAttrMethodParameters(io io.Reader, pool *ConstantPool) (*AttrMethodParameters, error) {
	var parametersCount uint8
	if err := binary.Read(io, binary.BigEndian, &parametersCount); err != nil {
		return nil, err
	}
	mp := &AttrMethodParameters{}
	for i := 0; i < int(parametersCount); i++ {
		p, err := NewParameter(io, pool)
		if err != nil {
			return nil, err
		}
		mp.Parameters = append(mp.Parameters, p)
	}
	return mp, nil
}

func (c *AttrMethodParameters) Name() string {
	return MethodParameters
}

func (c *AttrMethodParameters) String() string {
	return fmt.Sprintf("AttrMethodParameters ->  %v", *c)
}