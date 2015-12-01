package slipstream

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func BenchmarkBasicReader(b *testing.B) {
	html := "<html><head></head><body></body></html>"

	r := strings.NewReader(html)
	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		r.Seek(0, 0)
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSlipBefore(b *testing.B) {
	html := "<html><head></head><body></body></html>"

	var ins = []byte("<script></script>")
	var key = []byte("</body>")

	r := Slip(ins, Before(key), 0)(strings.NewReader(html))

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		r.(*Slipstream).Source.(*strings.Reader).Seek(0, 0)
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBasicReader-4  10000000               164 ns/op               0 B/op          0 allocs/op
// BenchmarkSlipBefore-4   10000000               177 ns/op               0 B/op          0 allocs/op
