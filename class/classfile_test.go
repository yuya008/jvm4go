package class

import (
	"testing"
	"os"
	"log"
	"io/ioutil"
	"path"
	"bytes"
)

var (
	testByteCode []byte
)

const (
	testClassFileName = "Test.class"
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
	t.Log(classFile)
}
