package class

import (
	"testing"
	"os"
	"log"
	"io/ioutil"
	"path"
	"bytes"
	"fmt"
)

var (
	testByteCode []byte
)

const (
	testClassFileName = "ArrayList.class"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if testByteCode, err = ioutil.ReadFile(path.Join(pwd, testClassFileName)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func TestNewClassFile(t *testing.T) {
	classFile, err := NewClassFile(bytes.NewReader(testByteCode))
	if err != nil {
		t.Error(err)
	}
	if classFile.Magic != 0XCAFEBABE {
		t.Error("format error")
	}
	showClassFile(classFile)
}

func showClassFile(cf *ClassFile) {
	fmt.Printf("Magic:%d\n", cf.Magic)
	fmt.Printf("Major:%d\n", cf.Major)
	fmt.Printf("Minor:%d\n", cf.Minor)

	for i := 0; i < cf.ConstantPool.Length(); i++ {
		constant, err := cf.ConstantPool.Get(uint16(i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		if constant == nil {
			fmt.Printf("ConstantPool[%d] %s\n", i, "<nil>")
		} else {
			fmt.Printf("ConstantPool[%d] %s\n", i, constant)
		}
	}
	fmt.Printf("AccessFlags:%s\n", showAccessFlags(cf.AccessFlags))
	fmt.Printf("ThisClass:%s\n", cf.ThisClass)
	fmt.Printf("SuperClass:%s\n", cf.SuperClass)

	for i, inter := range cf.Interfaces {
		fmt.Printf("Interface[%d]:%s\n", i, inter)
	}
	for i, field := range cf.Fields {
		fmt.Printf(`Fields[%d]:{
   AccessFlags: %s
   Name: %s
   Descriptor: %s
   Attrs: %v
}



`, i, showFieldAccessFlags(field.AccessFlags), field.Name, field.Descriptor, field.Attrs)
	}
	for i, m := range cf.Methods {
		fmt.Printf("method[%d] %s\n\n\n\n", i, m)
	}
	for i, attr := range cf.Attrs {
		fmt.Printf("Attr[%d] %s\n", i, attr)
	}
}

func showAccessFlags(accessFlags uint16) string {
	var s string
	if accessFlags & ACCPUBLIC > 0 {
		s += " ACC_PUBLIC "
	}
	if accessFlags & ACCFINAL > 0 {
		s += " ACC_FINAL "
	}
	if accessFlags & ACCSUPER > 0 {
		s += " ACC_SUPER "
	}
	if accessFlags & ACCINTERFACE > 0 {
		s += " ACC_INTERFACE "
	}
	if accessFlags & ACCABSTRACT > 0 {
		s += " ACC_ABSTRACT "
	}
	if accessFlags & ACCSYNTHETIC > 0 {
		s += " ACC_SYNTHETIC "
	}
	if accessFlags & ACCANNOTATION > 0 {
		s += " ACC_ANNOTATION "
	}
	if accessFlags & ACCENUM > 0 {
		s += " ACC_ENUM "
	}
	return s
}

func showFieldAccessFlags(accessFlags uint16) string {
	var s string
	if accessFlags & FieldAccPublic > 0 {
		s += " FieldAccPublic "
	}
	if accessFlags & FieldAccPrivate > 0 {
		s += " FieldAccPrivate "
	}
	if accessFlags & FieldAccProtected > 0 {
		s += " FieldAccProtected "
	}
	if accessFlags & FieldAccStatic > 0 {
		s += " FieldAccStatic "
	}
	if accessFlags & FieldAccFinal > 0 {
		s += " FieldAccFinal "
	}
	if accessFlags & FieldAccVolatile > 0 {
		s += " FieldAccVolatile "
	}
	if accessFlags & FieldAccTransient > 0 {
		s += " FieldAccTransient "
	}
	if accessFlags & FieldAccSynthetic > 0 {
		s += " FieldAccSynthetic "
	}
	if accessFlags & FieldAccEnum > 0 {
		s += " FieldAccEnum "
	}
	return s
}