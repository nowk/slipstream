package slipstream

import (
	"io"

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

	// flush truncated
	if lent := len(s.trunc); lent > 0 {
		n := lenp
		if lent < n {
			n = lent
		}
		for ; writ < n; writ++ {
			p[writ] = s.trunc[writ]
		}
		if s.trunc = s.trunc[writ:]; len(s.trunc) > 0 {
			return writ, nil
		}
	}

	if len(s.buf) == 0 && len(s.trunc) == 0 && s.eof {
		return writ, io.EOF
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

	// check occurrence count
	if s.max > 0 && s.count >= s.max {
		// move the buf to trunc and keep reading
		s.trunc = s.buf
		s.buf = s.buf[:0]

		return writ, nil
	}

	// slip the insert into the buf if applicable
	var out []byte
	var ok bool
	out, s.buf, ok = s.slipFunc(s.ins, s.buf)
	if ok {
		s.count++
	} else {
		// if no match, nothing read and we are EOF, we are done, lets write
		// everything to out
		if n == 0 && s.eof {
			out = append(out, s.buf...)

			s.buf = s.buf[:0]
		}
	}
	if n := lenp - writ; len(out) > n {
		s.trunc = out[n:] // set truncated

		out = out[:n]
	}

	// write out to p
	i := 0
	n = len(out) + writ
	for ; writ < n; writ++ {
		p[writ] = out[i]

		i++
	}

	return writ, nil
}

// SlipFunc is the func signature for inserting a bytes to a bytes. It takes the
// insert value and the source, respectively, as arguments.
//
// This func returns the bytes to write out to the Writer, bytes to be saved
// to buffer to be used in the next Read cycle and whether an insert occured or
// not.
type SlipFunc func([]byte, []byte) ([]byte, []byte, bool)

func Before(key []byte) SlipFunc {
	return func(ins, src []byte) ([]byte, []byte, bool) {
		i, m := bytematch.Compare(src, key)
		if m.Partial() {
			return src[:i], src[i:], false
		}

		if m.Exact() {
			src = append(src, ins...) // grow the slice by ins length

			// insert value, via copy to shift, and replace
			n := 0
			m := len(ins)
			for ; n < m; n++ {
				o := i + n

				copy(src[o+1:], src[o:])
				src[o] = ins[n]
			}

			// offset
			n = i + m + 1

			return src[:n], src[n:], true
		}

		return src, src[:0], false
	}
}
