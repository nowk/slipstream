# bytematch

[![Build Status](https://travis-ci.org/nowk/bytematch.svg?branch=master)](https://travis-ci.org/nowk/bytematch)
[![GoDoc](https://godoc.org/github.com/nowk/bytematch?status.svg)](http://godoc.org/github.com/nowk/bytematch)

Simple bytes comparison with match levels.

## Install

    go get github.com/nowk/bytematch


## Usage

Exact matches match for any exact occurrence in the source value.

    var a = []byte("Hello World!")
    var b = []byte("ello")

    n, m := bytematch.Compare(a, b)
    // m.Exact() => true
    // n         => 1


Partial matches match when a partial consecutive match at the end of the source 
value.

    var a = []byte("Hello Wor")
    var b = []byte("World!")

    n, m := bytematch.Compare(a, b)
    // m.Partial() => true
    // n           => 8


### License

MIT
