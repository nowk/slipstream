# slipstream

[![Build Status](https://travis-ci.org/nowk/slipstream.svg?branch=master)](https://travis-ci.org/nowk/slipstream)
[![GoDoc](https://godoc.org/github.com/nowk/slipstream?status.svg)](http://godoc.org/github.com/nowk/slipstream)

Slip values into your io.Reader stream


## Install

    go get github.com/nowk/slipstream


## Usage

    html := strings.NewReader("<html><head></head><body></body></html>")

    var (
        ins = []byte("<script></script>")
        key = []byte("</body>")
    )

    slip := slipstream.Slip(ins, slipstream.Before(key), 0)

    w := bytes.NewBuffer(nil)

    _, err := io.Copy(w, slip(html))
    if err != nil {
        // handle
    }


### License

MIT
