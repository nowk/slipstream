package slipstream

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestBeforeRead(t *testing.T) {
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

func TestBefore(t *testing.T) {
	var cases = []struct {
		source   string
		ins, key string
		exp      string
	}{
		{"123", "b", "c", "123"},

		{"ac", "b", "c", "abc"},
		{"acc", "b", "c", "abcbc"},
		{"accdefgc", "b", "c", "abcbcdefgbc"},

		{"Hello !", "World", "!", "Hello World!"},
		{"World!", "Hello ", "World", "Hello World!"},
	}

	for _, v := range cases {
		r := strings.NewReader(v.source)

		var exp = v.exp
		var ins = []byte(v.ins)
		var key = []byte(v.key)

		slip := Slip(ins, Before(key), 0)

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

func TestOccurrence(t *testing.T) {
	var cases = []struct {
		source   string
		ins, key string
		n        int
		exp      string
	}{
		{"accdefgc", "b", "c", 0, "abcbcdefgbc"},
		{"accdefgc", "b", "c", 2, "abcbcdefgc"},
		{"accdefgc", "b", "c", 1, "abccdefgc"},
	}

	for _, v := range cases {
		r := strings.NewReader(v.source)

		var exp = v.exp
		var ins = []byte(v.ins)
		var key = []byte(v.key)

		slip := Slip(ins, Before(key), v.n)

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
