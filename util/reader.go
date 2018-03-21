package util

import (
	"encoding/binary"
)

// BinaryReader reads primitive data types as binary values.
type BinaryReader struct {
	buffer []byte
	pos    int
}

// NewBinaryReader returns a new instance of the BinaryReader class based on the specified byte array.
func NewBinaryReader(packet []byte) *BinaryReader {
	return &BinaryReader{
		buffer: packet,
		pos:    0,
	}
}

// Pos returns the current position in the buffer
func (reader *BinaryReader) Pos() int {
	return reader.pos
}

func (reader *BinaryReader) Read(count int) []byte {
	b := reader.buffer[reader.pos : reader.pos+count]
	reader.pos += count
	return b
}

// ReadUInt8 reads the next uint8 from the current stream and advances the position of the stream by one byte.
func (reader *BinaryReader) ReadUInt8() uint8 {
	return reader.Read(1)[0]
}

// ReadUInt16 reads the next uint16 from the current stream and advances the position of the stream by two bytes.
func (reader *BinaryReader) ReadUInt16() uint16 {
	return binary.LittleEndian.Uint16(reader.Read(2))
}

// ReadUInt32 reads the next uint32 from the current stream and advances the position of the stream by four bytes.
func (reader *BinaryReader) ReadUInt32() uint32 {
	return binary.LittleEndian.Uint32(reader.Read(4))
}

// ReadUInt64 reads the next uint64 from the current stream and advances the position of the stream by eight bytes.
func (reader *BinaryReader) ReadUInt64() uint64 {
	return binary.LittleEndian.Uint64(reader.Read(8))
}

// ReadInt8 reads the next int8 from the current stream and advances the position of the stream by one byte.
func (reader *BinaryReader) ReadInt8() int8 {
	return int8(reader.ReadUInt8())
}

// ReadInt16 reads the next int16 from the current stream and advances the position of the stream by two bytes.
func (reader *BinaryReader) ReadInt16() int16 {
	return int16(reader.ReadUInt16())
}

// ReadInt32 reads the next int32 from the current stream and advances the position of the stream by four bytes.
func (reader *BinaryReader) ReadInt32() int32 {
	return int32(reader.ReadUInt32())
}

// ReadInt64 reads the next int64 from the current stream and advances the position of the stream by eight bytes.
func (reader *BinaryReader) ReadInt64() int64 {
	return int64(reader.ReadUInt64())
}

// ReadBool reads a bool from the current stream and advances the position of the stream by one byte.
func (reader *BinaryReader) ReadBool() bool {
	b := reader.ReadUInt8()
	return b != 0
}

// ReadCString reads a null terminated string from the current stream and advanced the position by the length of the string.
func (reader *BinaryReader) ReadCString() string {
	start := reader.pos
	for {
		if reader.buffer[reader.pos] == 0 {
			reader.pos++
			break
		}
		reader.pos++
	}
	return string(reader.buffer[start : reader.pos-1])
}

// More ...
func (reader *BinaryReader) More() bool {
	return reader.pos < len(reader.buffer)
}
