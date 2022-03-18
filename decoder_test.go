package linbuf

import "testing"

func FuzzDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte, a int16, b float64, d []byte) {
		t.Run("encode and decode", func(t *testing.T) {
			// encode and decode
			ec := NewEncoder().
				Int16(a).
				Float64(b).
				Bytes(data)
			defer ec.Destroy()
			dc := NewDecoder(ec.FinalizeBytes())
			var a2 int16
			var b2 float64
			dc.Int16(&a2).Float64(&b2)
			if a != a2 || b != b2 || dc.Finalize() != nil {
				t.Errorf("%v != %v", a, a2)
			}

			// random data decode
			NewDecoder(data).Float64(&b2).Finalize()
			NewDecoder(data).Int16(&a2).Finalize()
			NewDecoder(data).Bytes(&d)
			var s string
			NewDecoder(data).String(&s)
		})
	})
}
