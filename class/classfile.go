package class

import (
	"encoding/binary"
	"errors"
	"io"
	"go/constant"
)

type ClassFile struct {
	// 魔数
	Magic uint32
	// class文件主版本号
	Major uint16
	// class文件次版本号
	Minor uint16
	// 常量池
	ConstantPool *ConstantPool
	// 类访问标志
	AccessFlags uint16
	// 当前类索引
	ThisClass *ConstClass
	// 超类索引
	SuperClass *ConstClass
	// 接口索引表
	Interfaces []*ConstClass
	// 字段表
	Fields []*Field
	// 方法表
	Methods []*Method
	// 属性表
	Attrs []Attr
	reader io.Reader
}

const (
	ClassFileMagic = 0xcafebabe
)

const (
	ACCPUBLIC = 0x0001
	ACCFINAL  = 0x0010
	ACCSUPER  = 0x0020
	ACCINTERFACE = 0x0200
	ACCABSTRACT = 0x0400
	ACCSYNTHETIC = 0x1000
	ACCANNOTATION = 0x2000
	ACCENUM = 0x4000
)

var (
	ClassFileFormatError = errors.New("class file format invalid")
)

func NewClassFile(reader io.Reader) (*ClassFile, error) {
	var err error
	classFile := &ClassFile{
		reader: reader,
	}
	if err = classFile.readAndCheckMagic(); err != nil {
		return nil, err
	}
	if err = classFile.readAndCheckClassFileVersion(); err != nil {
		return nil, err
	}
	if classFile.ConstantPool, err = NewConstantPool(reader); err != nil {
		return nil, err
	}
	if err = classFile.readAccessFlags(); err != nil {
		return nil, err
	}
	if err = classFile.readThisClass(); err != nil {
		return nil, err
	}
	if err = classFile.readSuperClass(); err != nil {
		return nil, err
	}
	if err = classFile.readInterfaces(); err != nil {
		return nil, err
	}
	if err = classFile.readFields(); err != nil {
		return nil, err
	}
	return classFile, nil
}

func (classfile *ClassFile) readAndCheckMagic() error {
	if err := binary.Read(classfile.reader, binary.BigEndian, &classfile.Magic); err != nil {
		return err
	}
	if classfile.Magic != ClassFileMagic {
		return errors.New("class file 'magic' invalid")
	}
	return nil
}

func (classfile *ClassFile) readAndCheckClassFileVersion() error {
	if err := binary.Read(classfile.reader, binary.BigEndian, &classfile.Minor); err != nil {
		return err
	}
	if err := binary.Read(classfile.reader, binary.BigEndian, &classfile.Major); err != nil {
		return err
	}
	switch classfile.Major {
	case 45:
	case 46, 47, 48, 49, 50, 51, 52:
	default:
		return errors.New("unknow class version")
	}
	return nil
}

func (classfile *ClassFile) readAccessFlags() error {
	if err := binary.Read(classfile.reader, binary.BigEndian, &classfile.AccessFlags); err != nil {
		return err
	}
	return nil
}

func (classfile *ClassFile) readThisClass() error {
	var thisClassIndex uint16
	if err := binary.Read(classfile.reader, binary.BigEndian, &thisClassIndex); err != nil {
		return err
	}
	var err error
	if classfile.ThisClass, err = classfile.ConstantPool.GetClass(thisClassIndex); err != nil {
		return err
	}
	return nil
}

func (classfile *ClassFile) readSuperClass() error {
	var superClassIndex uint16
	var err error
	if err := binary.Read(classfile.reader, binary.BigEndian, &superClassIndex); err != nil {
		return err
	}
	if classfile.SuperClass, err = classfile.ConstantPool.GetClass(superClassIndex); err != nil {
		return err
	}
	return nil
}

func (classfile *ClassFile) readInterfaces() error {
	var interfaceCount uint16
	if err := binary.Read(classfile.reader, binary.BigEndian, &interfaceCount); err != nil {
		return err
	}
	var index uint16
	for i := 0; i < int(interfaceCount); i++ {
		if err := binary.Read(classfile.reader, binary.BigEndian, &index); err != nil {
			return err
		}
		class, err := classfile.ConstantPool.GetClass(index)
		if err != nil {
			return err
		}
		classfile.Interfaces = append(classfile.Interfaces, class)
	}
	return nil
}

func (classfile *ClassFile) readFields() error {
	var fieldsCount uint16
	if err := binary.Read(classfile.reader, binary.BigEndian, &fieldsCount); err != nil {
		return err
	}
	for i := 0; i < fieldsCount; i++ {
		field, err := classfile.readField()
		if err != nil {
			return err
		}
		classfile.Fields = append(classfile.Fields, field)
	}
	return nil
}

func (classfile *ClassFile) readField() (*Field, error) {
	field := &Field{}
	if err := binary.Read(classfile.reader, binary.BigEndian, field.AccessFlags); err != nil {
		return nil, err
	}
	var nameIndex, descriptorIndex uint16
	if err := binary.Read(classfile.reader, binary.BigEndian, &nameIndex); err != nil {
		return nil, err
	}
	if err := binary.Read(classfile.reader, binary.BigEndian, &descriptorIndex); err != nil {
		return nil, err
	}
	var err error
	if field.Name, err = classfile.ConstantPool.GetUTF8String(nameIndex); err != nil {
		return nil, err
	}
	if field.Descriptor, err = classfile.ConstantPool.GetUTF8String(descriptorIndex); err != nil {
		return nil, err
	}
	var attrCount uint16
	if err := binary.Read(classfile.reader, binary.BigEndian, &attrCount); err != nil {
		return nil, err
	}
	field.Attrs = make([]Attr, attrCount)
	for i := 0; i < attrCount; i++ {
		attr, err := ReadAttr(classfile.reader, classfile.ConstantPool)
		if err != nil {
			return nil, err
		}
		field.Attrs[i] = attr
	}
	return field, nil
}
