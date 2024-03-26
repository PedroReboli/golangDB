package godb

import (
	"encoding/binary"
)

type ByteReader struct{
	pos uint64
	data []byte
}
func NewByteReader(data []byte) ByteReader{
	return ByteReader{
		pos:0,
		data: data,
	}
}

func (br *ByteReader) GetUint32() uint32{
	v := binary.BigEndian.Uint32(br.data[br.pos:br.pos+4])
	br.pos += 4
	return v
}

func (br *ByteReader) GetUint16() uint16{
	
	v := binary.BigEndian.Uint16(br.data[br.pos:br.pos+2])
	br.pos += 2
	return v
}

func (br *ByteReader) GetUint64() uint64{
	
	v := binary.BigEndian.Uint64(br.data[br.pos:br.pos+8])
	br.pos += 8
	return v
}

func (br *ByteReader) GetByte() byte{
	v := br.data[br.pos]
	br.pos += 1
	return v
}

func (br *ByteReader) GetBytes(size int16) []byte{
	v := br.data[br.pos:br.pos+uint64(size)]
	br.pos += uint64(size)
	return v
}

func (br *ByteReader) GetString(size int16) string{
	println("Get string size", size)
	v := br.data[br.pos:br.pos+uint64(size)]
	br.pos += uint64(size)
	return string(v)
}

func GetUint32At(data []byte,pos uint64) uint32{
	v := binary.BigEndian.Uint32(data[pos:pos+4])
	return v
}

func GetUint16At(data []byte, pos uint64) uint16{
	
	v := binary.BigEndian.Uint16(data[pos:pos+2])
	
	return v
}

func GetUint64At(data []byte, pos uint64) uint64{
	
	
	return binary.BigEndian.Uint64(data[pos:pos+8])
}

func GetByteAt(data []byte, pos uint64) byte{
	
	
	return data[pos]
}

func GetBytesAt(data []byte, pos uint64, size int16) []byte{
	return data[pos:pos+uint64(size)]
}

func GetStringAt(data []byte, pos uint64, size int16) string{
	println("string size", size)
	v := data[pos:pos+uint64(size)]
	pos += uint64(size)
	return string(v)
}

