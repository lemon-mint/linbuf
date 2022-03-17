package linbuf

import (
	"encoding/binary"
	"io"
	"math"
	"sync"
	"unsafe"
)

type byteReader struct {
	Data []byte
	Pos  int
}

func (r *byteReader) ReadByte() (byte, error) {
	if r.Pos >= len(r.Data) {
		return 0, io.EOF
	}
	v := r.Data[r.Pos]
	r.Pos++
	return v, nil
}

func (r *byteReader) Read(p []byte) (int, error) {
	n := copy(p, r.Data[r.Pos:])
	r.Pos += n
	if n == 0 && len(p) > 0 {
		return n, io.EOF
	}
	return n, nil
}

func (r *byteReader) Reset() {
	r.Pos = 0
	r.Data = nil
}

type Decoder struct {
	b byteReader
	e error
}

func resetDecoder(d *Decoder) {
	d.b.Reset()
	d.e = nil
}

var decoderPool = sync.Pool{
	New: func() any {
		return &Decoder{}
	},
}

func putDecoder(d *Decoder) {
	resetDecoder(d)
	decoderPool.Put(d)
}

func getDecoder() *Decoder {
	return decoderPool.Get().(*Decoder)
}

func NewDecoder(p []byte) *Decoder {
	b := getDecoder()
	b.b.Data = p
	return b
}

func (b *Decoder) Uint8(p *uint8) *Decoder {
	if b.e != nil {
		return b
	}
	v, err := b.b.ReadByte()
	if err != nil {
		b.e = err
		return b
	}
	*p = uint8(v)
	return b
}

func (b *Decoder) Uint16(p *uint16) *Decoder {
	if b.e != nil {
		return b
	}
	var v [2]byte
	_, err := b.b.Read(v[:])
	if err != nil {
		b.e = err
		return b
	}
	*p = uint16(v[0]) | uint16(v[1])<<8
	return b
}

func (b *Decoder) Uint32(p *uint32) *Decoder {
	if b.e != nil {
		return b
	}
	var v [4]byte
	_, err := b.b.Read(v[:])
	if err != nil {
		b.e = err
		return b
	}
	*p = uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24
	return b
}

func (b *Decoder) Uint64(p *uint64) *Decoder {
	if b.e != nil {
		return b
	}
	var v [8]byte
	_, err := b.b.Read(v[:])
	if err != nil {
		b.e = err
		return b
	}
	*p = uint64(v[0]) | uint64(v[1])<<8 | uint64(v[2])<<16 | uint64(v[3])<<24 |
		uint64(v[4])<<32 | uint64(v[5])<<40 | uint64(v[6])<<48 | uint64(v[7])<<56
	return b
}

func (b *Decoder) Int8(p *int8) *Decoder {
	return b.Uint8((*uint8)(unsafe.Pointer(p)))
}

func (b *Decoder) Int16(p *int16) *Decoder {
	return b.Uint16((*uint16)(unsafe.Pointer(p)))
}

func (b *Decoder) Int32(p *int32) *Decoder {
	return b.Uint32((*uint32)(unsafe.Pointer(p)))
}

func (b *Decoder) Int64(p *int64) *Decoder {
	return b.Uint64((*uint64)(unsafe.Pointer(p)))
}

func (b *Decoder) Float32(p *float32) *Decoder {
	var v uint32
	b.Uint32(&v)
	*p = math.Float32frombits(v)
	return b
}

func (b *Decoder) Float64(p *float64) *Decoder {
	var v uint64
	b.Uint64(&v)
	*p = math.Float64frombits(v)
	return b
}

func (b *Decoder) VarUint64(p *uint64) *Decoder {
	if b.e != nil {
		return b
	}
	*p, b.e = binary.ReadUvarint(&b.b)
	return b
}

func (b *Decoder) VarUint32(p *uint32) *Decoder {
	if b.e != nil {
		return b
	}
	var v uint64
	v, b.e = binary.ReadUvarint(&b.b)
	*p = uint32(v)
	return b
}

func (b *Decoder) VarUint16(p *uint16) *Decoder {
	if b.e != nil {
		return b
	}
	var v uint64
	v, b.e = binary.ReadUvarint(&b.b)
	*p = uint16(v)
	return b
}

func (b *Decoder) VarUint8(p *uint8) *Decoder {
	if b.e != nil {
		return b
	}
	var v uint64
	v, b.e = binary.ReadUvarint(&b.b)
	*p = uint8(v)
	return b
}

func (b *Decoder) VarInt64(p *int64) *Decoder {
	if b.e != nil {
		return b
	}
	*p, b.e = binary.ReadVarint(&b.b)
	return b
}

func (b *Decoder) VarInt32(p *int32) *Decoder {
	if b.e != nil {
		return b
	}
	var v int64
	v, b.e = binary.ReadVarint(&b.b)
	*p = int32(v)
	return b
}

func (b *Decoder) VarInt16(p *int16) *Decoder {
	if b.e != nil {
		return b
	}
	var v int64
	v, b.e = binary.ReadVarint(&b.b)
	*p = int16(v)
	return b
}

func (b *Decoder) Bytes(p *[]byte) *Decoder {
	var n uint64
	b.VarUint64(&n)
	if b.e != nil {
		return b
	}
	if n > uint64(len(b.b.Data)-b.b.Pos) {
		b.e = io.ErrUnexpectedEOF
		return b
	}
	*p = b.b.Data[b.b.Pos : b.b.Pos+int(n)]
	b.b.Pos += int(n)
	return b
}

func (b *Decoder) String(p *string) *Decoder {
	var b2 []byte
	b.Bytes(&b2)
	*p = string(b2)
	return b
}

func (b *Decoder) Finalize() error {
	err := b.e
	putDecoder(b)
	return err
}
