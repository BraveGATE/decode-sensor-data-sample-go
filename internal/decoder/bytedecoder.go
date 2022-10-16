package decoder

import (
	"bytes"
	"encoding/binary"
)

func Int8FromByte(buf []byte) int8 {
	var i int8
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &i); err != nil {
		return 0
	}
	return i
}

func Int16FromByte(buf []byte) int16 {
	var i int16
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &i); err != nil {
		return 0
	}
	return i
}

func Int32FromByte(buf []byte) int32 {
	var i int32
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &i); err != nil {
		return 0
	}
	return i
}

func Float32FromByte(buf []byte) float32 {
	var i float32
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &i); err != nil {
		return 0
	}
	return i
}
