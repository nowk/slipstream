package bytematch

import (
	"bytes"
)

type MatchLevel int

const (
	_ MatchLevel = iota

	Exact
	Partial
	Not
)

func (m MatchLevel) Exact() bool {
	return m == Exact
}

func (m MatchLevel) Partial() bool {
	return m == Partial
}

func (m MatchLevel) Not() bool {
	return m == Not
}

func Compare(src, v []byte) (int, MatchLevel) {
	if i := bytes.Index(src, v); i != -1 {
		return i, Exact
	}

	lensrc := len(src)

	z := len(v)
	for ; z > 0; z-- {
		i := bytes.Index(src, v[:z])
		if i == -1 {
			continue
		}
		if i+z == lensrc {
			return i, Partial
		}
	}

	return -1, Not
}
