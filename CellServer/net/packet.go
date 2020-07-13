package net

import (
	"bytes"
	"encoding/base64"
)

func CreateReadPacket(b []byte) *Packet {
	return &Packet{
		buffer: bytes.NewBuffer(b),
	}
}

func CreateWritePacket() *Packet {
	return &Packet{
		buffer: bytes.NewBufferString(""),
	}
}

type Packet struct {
	buffer *bytes.Buffer
}

func (p *Packet) BytesBuffer() *bytes.Buffer {
	return p.buffer
}

func (p *Packet) BytesSlice() []byte {
	return p.buffer.Bytes()
}

func (p *Packet) CheckFlag() int8 {
	return int8(p.buffer.Bytes()[0])
}

func (p *Packet) ReadString() string {
	var buf [2]byte
	p.buffer.Read(buf[:])
	isize := BytesToInt16(buf[:])
	strbuf := make([]byte, isize)
	p.buffer.Read(strbuf)
	decodeBytes, err := base64.StdEncoding.DecodeString(string(strbuf))
	if err != nil {
		return ""
	}
	return string(decodeBytes)
}

func (p *Packet) ReadByte() byte {
	bit, _ := p.buffer.ReadByte()
	return bit
}

func (p *Packet) ReadBool() bool {
	i, _ := p.buffer.ReadByte()
	if i == 0 {
		return false
	} else {
		return true
	}
}

func (p *Packet) ReadBytes() []byte {
	var buf [2]byte
	p.buffer.Read(buf[:])
	isize := BytesToInt16(buf[:])
	strbuf := make([]byte, isize)
	p.buffer.Read(strbuf)
	return strbuf
}

func (p *Packet) ReadInt8() int8 {
	i, _ := p.buffer.ReadByte()
	return int8(i)
}

func (p *Packet) ReadInt16() int16 {
	var buf [2]byte
	p.buffer.Read(buf[:])
	return BytesToInt16(buf[:])
}

func (p *Packet) ReadInt32() int32 {
	var buf [4]byte
	p.buffer.Read(buf[:])
	return BytesToInt32(buf[:])
}

func (p *Packet) ReadInt64() int64 {
	var buf [8]byte
	p.buffer.Read(buf[:])
	return BytesToInt64(buf[:])
}

func (p *Packet) ReadFloat32() float32 {
	var buf [4]byte
	p.buffer.Read(buf[:])
	return BytesToFloat32(buf[:])
}

func (p *Packet) ReadFloat64() float64 {
	var buf [8]byte
	p.buffer.Read(buf[:])
	return BytesToFloat64(buf[:])
}

func (p *Packet) WriteString(s string) {
	bstr := base64.StdEncoding.EncodeToString([]byte(s))
	size := int16(len(bstr))
	p.buffer.Write(Int16ToBytes(size))
	p.buffer.Write([]byte(bstr))
}

func (p *Packet) WriteBytes(b []byte) {
	size := int16(len(b))
	p.buffer.Write(Int16ToBytes(size))
	p.buffer.Write(b)
}

func (p *Packet) WriteBuffer(b []byte) {
	p.buffer.Write(b)
}

func (p *Packet) WriteByte(b byte) {
	p.buffer.Write(Int8ToBytes(int8(b)))
}

func (p *Packet) WriteBool(b bool) {
	if !b {
		p.buffer.Write(Int8ToBytes(0))
	} else {
		p.buffer.Write(Int8ToBytes(1))
	}
}

func (p *Packet) WriteInt8(i int8) {
	p.buffer.Write(Int8ToBytes(i))
}

func (p *Packet) WriteInt16(i int16) {
	p.buffer.Write(Int16ToBytes(i))
}

func (p *Packet) WriteInt32(i int32) {
	p.buffer.Write(Int32ToBytes(i))
}

func (p *Packet) WriteInt64(i int64) {
	p.buffer.Write(Int64ToBytes(i))
}

func (p *Packet) WriteFloat32(f float32) {
	p.buffer.Write(Float32ToBytes(f))
}

func (p *Packet) WriteFloat64(f float64) {
	p.buffer.Write(Float64ToBytes(f))
}
