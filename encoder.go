package linbuf

import (
	"encoding/binary"
	"math"

	"github.com/valyala/bytebufferpool"
)

type Encoder struct {
	b *bytebufferpool.ByteBuffer
}

func NewEncoder() Encoder {
	return Encoder{
		b: bytebufferpool.Get(),
	}
}

func (e Encoder) Reset() {
	e.b.Reset()
}

func (e Encoder) Finalize() *bytebufferpool.ByteBuffer {
	return e.b
}

func (e Encoder) FinalizeBytes() []byte {
	return e.b.Bytes()
}

func (e Encoder) Destroy() {
	bytebufferpool.Put(e.b)
}

func (e Encoder) Uint8(v uint8) Encoder {
	e.b.WriteByte(v)
	return e
}

func (e Encoder) Uint16(v uint16) Encoder {
	var b [2]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	e.b.Write(b[:])
	return e
}

func (e Encoder) Uint32(v uint32) Encoder {
	var b [4]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	e.b.Write(b[:])
	return e
}

func (e Encoder) Uint64(v uint64) Encoder {
	var b [8]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	e.b.Write(b[:])
	return e
}

func (e Encoder) Int8(v int8) Encoder {
	return e.Uint8(uint8(v))
}

func (e Encoder) Int16(v int16) Encoder {
	return e.Uint16(uint16(v))
}

func (e Encoder) Int32(v int32) Encoder {
	return e.Uint32(uint32(v))
}

func (e Encoder) Int64(v int64) Encoder {
	return e.Uint64(uint64(v))
}

func (e Encoder) Float32(v float32) Encoder {
	return e.Uint32(math.Float32bits(v))
}

func (e Encoder) Float64(v float64) Encoder {
	return e.Uint64(math.Float64bits(v))
}

func (e Encoder) VarUint64(v uint64) Encoder {
	var b [10]byte
	n := binary.PutUvarint(b[:], v)
	e.b.Write(b[:n])
	return e
}

func (e Encoder) VarUint32(v uint32) Encoder {
	var b [6]byte
	n := binary.PutUvarint(b[:], uint64(v))
	e.b.Write(b[:n])
	return e
}

func (e Encoder) VarUint16(v uint16) Encoder {
	var b [4]byte
	n := binary.PutUvarint(b[:], uint64(v))
	e.b.Write(b[:n])
	return e
}

func (e Encoder) VarInt64(v int64) Encoder {
	var b [10]byte
	n := binary.PutVarint(b[:], v)
	e.b.Write(b[:n])
	return e
}

func (e Encoder) VarInt32(v int32) Encoder {
	var b [6]byte
	n := binary.PutVarint(b[:], int64(v))
	e.b.Write(b[:n])
	return e
}

func (e Encoder) VarInt16(v int16) Encoder {
	var b [4]byte
	n := binary.PutVarint(b[:], int64(v))
	e.b.Write(b[:n])
	return e
}

func (e Encoder) String(s string) Encoder {
	e.VarUint64(uint64(len(s)))
	e.b.WriteString(s)
	return e
}

func (e Encoder) Bytes(b []byte) Encoder {
	e.VarUint64(uint64(len(b)))
	e.b.Write(b)
	return e
}

func (e Encoder) Boolean(b bool) Encoder {
	if b {
		e.Uint8(1)
	} else {
		e.Uint8(0)
	}
	return e
}
