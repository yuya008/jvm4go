package class

import (
	"io"
	"encoding/binary"
	"errors"
	"fmt"

	"go/constant"
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
}

func ReadAttr(reader io.Reader, pool *ConstantPool) (Attr, error) {
	var attrNameIndex uint16
	if err := binary.Read(reader, binary.BigEndian, &attrNameIndex); err != nil {
		return nil, err
	}
	attrName, err := pool.GetUTF8String(attrNameIndex)
	if err != nil {
		return err
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
	case BootstrapMethods:
	case Deprecated:
	case EnclosingMethod:
		return NewAttrEnclosingMethod(reader, pool)
	case Exceptions:
		return NewAttrExceptions(reader, pool)
	case InnerClasses:
		return NewAttrInnerClasses(reader, pool)
	case LineNumberTable:
		return NewAttrLineNumberTable(reader)
	case LocalVariableTable:
	case LocalVariableTypeTable:
	case MethodParameters:
	case RuntimeInvisibleAnnotations:
	case RuntimeInvisibleParameterAnnotations:
	case RuntimeInvisibleTypeAnnotations:
	case RuntimeVisibleAnnotations:
	case RuntimeVisibleParameterAnnotations:
	case RuntimeVisibleTypeAnnotations:
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
	for i := 0; i < exceptionLen; i++ {
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
	for i := 0; i < attrCount; i++ {
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
		s, err := NewAttrStackMapTable(io, pool)
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
	for i := 0; i < numberOfClasses; i++ {
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

type AttrSynthetic struct {}

func NewAttrSynthetic() (*AttrSynthetic, error) {
	return &AttrSynthetic{}, nil
}

func (c *AttrSynthetic) Name() string {
	return Synthetic
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

type AttrLocalVariableTypeTable struct {

}

type AttrDeprecated struct {

}

type AttrRuntimeVisibleAnnotations struct {

}

type AttrRuntimeInvisibleAnnotations struct {

}

type AttrRuntimeVisibleParameterAnnotations struct {

}

type AttrRuntimeInvisibleParameterAnnotations struct {

}

type AttrRuntimeVisibleTypeAnnotations struct {

}

type AttrRuntimeInvisibleTypeAnnotations struct {

}

type AttrAnnotationDefault struct {

}

type AttrBootstrapMethods struct {

}

type AttrMethodParameters struct {

}
