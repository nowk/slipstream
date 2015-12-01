package slipstream

import (
	"io"
	"log"

	"github.com/nowk/bytematch"
)

type Slipstream struct {
	// Source is the reader to slip values into
	Source io.Reader

	slipFunc SlipFunc
	ins      []byte
	max      int
	count    int

	// trunc are bytes that could not be written out due to the len of p
	trunc []byte

	// buf are partial matched bytes required for the next read comparison
	buf []byte

	// eof marks that the Source has reached EOF
	eof bool
}

var _ io.Reader = &Slipstream{}

func Slip(ins []byte, fn SlipFunc, n int) func(io.Reader) io.Reader {
	return func(r io.Reader) io.Reader {
		return &Slipstream{
			Source: r,

			slipFunc: fn,
			ins:      ins,
			max:      n,
		}
	}
}

func (s *Slipstream) Read(p []byte) (int, error) {
	lenp := len(p)
	writ := 0

	// log.Printf("buf %s, %d", s.buf, lenp)

	log.Printf("trunc %s, %d", s.trunc, len(s.trunc))
	if lent := len(s.trunc); lent > 0 {
		n := lenp
		if lent < n {
			n = lent
		}

		for ; writ < n; writ++ {
			p[writ] = s.trunc[writ]
		}

		s.trunc = s.trunc[writ:]

		if len(s.trunc) > 0 {
			return writ, nil
		}
	}

	n, err := s.Source.Read(p[writ:])
	if err != nil && err != io.EOF {
		return n, err
	}
	if err == io.EOF {
		s.eof = true
	}

	// alloc a buffer based on our initial read size
	if s.buf == nil {
		s.buf = make([]byte, 0, n)
	}
	if n > 0 {
		s.buf = append(s.buf, p[writ:writ+n]...)
	}

	// slip the insert into the buf if applicable
	out, buf := s.slipFunc(s.ins, s.buf)

	// save remaining buffer
	s.buf = buf

	log.Printf("append %s %d", s.buf, writ)

	var b []byte

	// calculate how much we can write to p
	rem := lenp - writ

	if len(out) > rem {
		b = out[:rem]

		// set truncated
		s.trunc = out[rem:]
	} else {
		b = out
	}

	// write out to p
	i := 0
	n = len(b) + writ
	for ; writ < n; writ++ {
		p[writ] = b[i]
		i++
	}
	if len(s.buf) == 0 && s.eof {
		return writ, io.EOF
	}

	return writ, nil
}

type SlipFunc func([]byte, []byte) (out []byte, buf []byte)

func Before(key []byte) SlipFunc {
	return func(ins, src []byte) ([]byte, []byte) {
		i, m := bytematch.Compare(src, key)
		if m.Partial() {
			return src[:i], src[i:]
		}

		if m.Exact() {
			src = append(src, ins...) // grow the slice by ins length

			// insert ins
			n := 0
			m := len(ins)
			for ; n < m; n++ {
				o := i + n

				copy(src[o+1:], src[o:])
				src[o] = ins[n]
			}

			n = i + m + 1

			return src[:n], src[n:]
		}

		return src, src[:0]
	}
}
