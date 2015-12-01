package slipstream

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	t.Skip()
	html := strings.NewReader("<html><body></body></html>")

	var ins = []byte("<script></script>")
	var key = []byte("</body>")

	slip := Slip(ins, Before(key), 0)

	r := slip(html)

	var cases = []struct {
		n   int
		exp string
	}{
		{15, "<html><body>"},
		{5, "<scri"},
		{5, "pt></"},
		{10, "script></b"},
		{3, "ody"},
		{15, "></html>"},
	}

	for _, v := range cases {
		b := make([]byte, v.n)

		var exp = v.exp

		n, err := r.Read(b)
		if err != nil {
			t.Fatal(err)
		}
		if got := string(b[:n]); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func TestSlipBeforeMatch(t *testing.T) {
	var cases = []struct {
		giv string
		exp string
	}{
		{"123", "123"},
		{"ac", "abc"},
		{"acc", "abcbc"},
		{"accdefgc", "abcbcdefgbc"},
	}

	for _, v := range cases {
		r := strings.NewReader(v.giv)

		var exp = v.exp

		var (
			ins = []byte("b")
			key = []byte("c")
		)

		slip := Slip(ins, Before(key), 1)

		buf := bytes.NewBuffer(nil)

		_, err := io.Copy(buf, slip(r))
		if err != nil {
			t.Fatal(err)
		}
		if got := buf.String(); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}
