package class

import (
	"errors"
	"io"
	"encoding/binary"
)

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
	if int(i) >= cp.Length() {
		return nil, errors.New("out of ConstantPool index")
	}
	return cp.pool[i], nil
}

func (cp *ConstantPool) GetUTF8String(i uint16) (*ConstUTF8, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstUTF8); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstUTF8) failed")
}

func (cp *ConstantPool) GetFloat(i uint16) (*ConstFloat, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstFloat); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstFloat) failed")
}

func (cp *ConstantPool) GetInteger(i uint16) (*ConstInteger, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstInteger); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstInteger) failed")
}

func (cp *ConstantPool) GetLong(i uint16) (*ConstLong, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstLong); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstLong) failed")
}

func (cp *ConstantPool) GetDouble(i uint16) (*ConstDouble, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstDouble); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstDouble) failed")
}

func (cp *ConstantPool) GetFieldRef(i uint16) (*ConstFieldRef, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstFieldRef); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstFieldRef) failed")
}

func (cp *ConstantPool) GetString(i uint16) (*ConstString, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstString); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstString) failed")
}

func (cp *ConstantPool) GetClass(i uint16) (*ConstClass, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstClass); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstClass) failed")
}

func (cp *ConstantPool) GetMethodRef(i uint16) (*ConstMethodRef, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstMethodRef); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstMethodRef) failed")
}

func (cp *ConstantPool) GetInterfaceMethodRef(i uint16) (*ConstInterfaceMethodRef, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstInterfaceMethodRef); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstInterfaceMethodRef) failed")
}

func (cp *ConstantPool) GetNameAndType(i uint16) (*ConstNameAndType, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstNameAndType); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstNameAndType) failed")
}

func (cp *ConstantPool) GetMethodHandle(i uint16) (*ConstMethodHandle, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstMethodHandle); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstMethodHandle) failed")
}

func (cp *ConstantPool) GetMethodType(i uint16) (*ConstMethodType, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
		return nil, err
	}
	if val, ok := constant.(*ConstMethodType); ok {
		return val, nil
	}
	return nil, errors.New("covert to constant.(*ConstMethodType) failed")
}

func (cp *ConstantPool) GetInvokeDynamic(i uint16) (*ConstInvokeDynamic, error) {
	constant, err := cp.Get(i)
	if err == nil {
		if constant == nil {
			return nil, nil
		}
	} else {
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
