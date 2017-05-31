package class

import (
	"io"
	"encoding/binary"
	"errors"
)

type StackMapFrame interface {
	FrameType() uint8
}

func NewStackMapFrame(io io.Reader, pool *ConstantPool) (StackMapFrame, error) {
	var frameType uint8
	if err := binary.Read(io, binary.BigEndian, &frameType); err != nil {
		return nil, err
	}
	switch {
	case frameType >= 0 && frameType <= 63:
		return &SameFrame{frameType:frameType}, nil
	case frameType >= 64 && frameType <= 127:
		return newSameLocals1StackItemFrame(frameType, io, pool)
	case frameType == 247:
		return newSameLocals1StackItemFrameExtended(frameType, io, pool)
	case frameType >= 248 && frameType <= 250:
		return newChopFrame(frameType, io)
	case frameType == 251:
		return newSameFrameExtended(frameType, io)
	case frameType >= 252 && frameType <= 254:
		return newAppendFrame(frameType, io, pool)
	case frameType == 255:
		return newFullFrame(frameType, io, pool)
	}
	return nil, errors.New("invalid frameType")
}

type SameFrame struct {
	frameType uint8
}

func (s *SameFrame) FrameType() uint8 {
	return s.frameType
}

type SameLocals1StackItemFrame struct {
	frameType uint8
	Stack []VerificationType
}

func newSameLocals1StackItemFrame(frameType uint8, io io.Reader, pool *ConstantPool) (*SameLocals1StackItemFrame, error) {
	vt, err := NewVerificationType(io, pool)
	if err != nil {
		return nil, err
	}
	return &SameLocals1StackItemFrame{
		frameType: frameType,
		Stack: []VerificationType{vt},
	}, nil
}

func (s *SameLocals1StackItemFrame) FrameType() uint8 {
	return s.frameType
}

type SameLocals1StackItemFrameExtended struct {
	frameType uint8
	OffsetDelta uint16
	Stack []VerificationType
}

func newSameLocals1StackItemFrameExtended(frameType uint8, io io.Reader, pool *ConstantPool) (*SameLocals1StackItemFrameExtended, error) {
	s := &SameLocals1StackItemFrameExtended{frameType:frameType}
	if err := binary.Read(io, binary.BigEndian, &s.OffsetDelta); err != nil {
		return nil, err
	}
	vt, err := NewVerificationType(io, pool)
	if err != nil {
		return nil, err
	}
	s.Stack = []VerificationType{vt}
	return s, nil
}

func (s *SameLocals1StackItemFrameExtended) FrameType() uint8 {
	return s.frameType
}

type ChopFrame struct {
	frameType uint8
	OffsetDelta uint16
}

func newChopFrame(frameType uint8, io io.Reader) (*ChopFrame, error) {
	c := &ChopFrame{frameType:frameType}
	if err := binary.Read(io, binary.BigEndian, &c.OffsetDelta); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ChopFrame) FrameType() uint8 {
	return s.frameType
}

type SameFrameExtended struct {
	frameType uint8
	OffsetDelta uint16
}

func newSameFrameExtended(frameType uint8, io io.Reader) (*SameFrameExtended, error) {
	c := &SameFrameExtended{frameType:frameType}
	if err := binary.Read(io, binary.BigEndian, &c.OffsetDelta); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *SameFrameExtended) FrameType() uint8 {
	return s.frameType
}

type AppendFrame struct {
	frameType uint8
	OffsetDelta uint16
	Locals []VerificationType
}

func newAppendFrame(frameType uint8, io io.Reader, pool *ConstantPool) (*AppendFrame, error) {
	a := &AppendFrame{frameType:frameType}
	if err := binary.Read(io, binary.BigEndian, &a.OffsetDelta); err != nil {
		return nil, err
	}
	numberOfLocals := frameType - uint8(251)
	if numberOfLocals <= 0 {
		return nil, errors.New("appendFrame number of locals == 0")
	}
	for i := 0; i < int(numberOfLocals); i++ {
		vt, err := NewVerificationType(io, pool)
		if err != nil {
			return nil, err
		}
		a.Locals = append(a.Locals, vt)
	}
	return a, nil
}

func (s *AppendFrame) FrameType() uint8 {
	return s.frameType
}

type FullFrame struct {
	frameType uint8
	OffsetDelta uint16
	Locals []VerificationType
	Stack []VerificationType
}

func newFullFrame(frameType uint8, io io.Reader, pool *ConstantPool) (*FullFrame, error) {
	f := FullFrame{frameType:frameType}
	if err := binary.Read(io, binary.BigEndian, &f.OffsetDelta); err != nil {
		return nil, err
	}
	var numberOfLocals, numberOfStackItems uint16
	if err := binary.Read(io, binary.BigEndian, &numberOfLocals); err != nil {
		return nil, err
	}
	for i := 0; i < len(numberOfLocals); i++ {
		vt, err := NewVerificationType(io, pool)
		if err != nil {
			return nil, err
		}
		f.Locals = append(f.Locals, vt)
	}
	if err := binary.Read(io, binary.BigEndian, &numberOfStackItems); err != nil {
		return nil, err
	}
	for i := 0; i < len(numberOfStackItems); i++ {
		vt, err := NewVerificationType(io, pool)
		if err != nil {
			return nil, err
		}
		f.Stack = append(f.Stack, vt)
	}
	return f, nil
}

func (s *FullFrame) FrameType() uint8 {
	return s.frameType
}