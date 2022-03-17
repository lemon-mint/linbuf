package linbuf

import (
	"testing"
)

func BenchmarkUint64(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i uint64
		for p.Next() {
			e := NewEncoder().
				Uint64(i)
			var v uint64
			NewDecoder(e.FinalizeBytes()).
				Uint64(&v).
				Finalize()
			if v != uint64(i) {
				b.Fatalf("expected %d, got %d", i, v)
			}
			e.Destroy()
			i++
		}
	})
}

func BenchmarkInt64(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i int64
		for p.Next() {
			e := NewEncoder().
				Int64(i)
			var v int64
			NewDecoder(e.FinalizeBytes()).
				Int64(&v).
				Finalize()
			if v != int64(i) {
				b.Fatalf("expected %d, got %d", i, v)
			}
			e.Destroy()
			i++
		}
	})
}

func BenchmarkVarUint64(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i uint64
		for p.Next() {
			e := NewEncoder().
				VarUint64(i)
			var v uint64
			NewDecoder(e.FinalizeBytes()).
				VarUint64(&v).
				Finalize()
			if v != uint64(i) {
				b.Fatalf("expected %d, got %d", i, v)
			}
			e.Destroy()
			i++
		}
	})
}

func BenchmarkVarInt64(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i int64
		for p.Next() {
			e := NewEncoder().
				VarInt64(i)
			var v int64
			NewDecoder(e.FinalizeBytes()).
				VarInt64(&v).
				Finalize()
			if v != int64(i) {
				b.Fatalf("expected %d, got %d", i, v)
			}
			e.Destroy()
			i++
		}
	})
}

func BenchmarkBytes(b *testing.B) {
	var data = []byte("Hello, World!")
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			e := NewEncoder().
				Bytes(data)
			var v []byte
			NewDecoder(e.FinalizeBytes()).
				Bytes(&v).
				Finalize()
			if string(v) != string(data) {
				b.Fatalf("expected %s, got %s", data, v)
			}
			e.Destroy()
		}
	})
}

func BenchmarkString(b *testing.B) {
	var data = "Hello, World!"
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			e := NewEncoder().
				String(data)
			var v string
			NewDecoder(e.FinalizeBytes()).
				String(&v).
				Finalize()
			if v != data {
				b.Fatalf("expected %s, got %s", data, v)
			}
			e.Destroy()
		}
	})
}
