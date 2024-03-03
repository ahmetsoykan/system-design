package hash

import (
	"bytes"
	"strings"
)

const chars string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// using base62 encoding for unique IDs
func Encode(nb uint64) string {
	base := chars
	var buf bytes.Buffer
	l := uint64(len(base))
	if nb/l != 0 {
		encode(nb/l, &buf, base)
	}
	buf.WriteByte(base[nb%l])
	return buf.String()
}

func encode(nb uint64, buf *bytes.Buffer, base string) {
	l := uint64(len(base))
	if nb/l != 0 {
		encode(nb/l, buf, base)
	}
	buf.WriteByte(base[nb%l])
}

// base62 decoding back to find initial input
func Decode(enc string) uint64 {
	base := chars
	var nb uint64
	lbase := len(base)
	le := len(enc)
	for i := 0; i < le; i++ {
		mult := 1
		for j := 0; j < le-i-1; j++ {
			mult *= lbase
		}
		nb += uint64(strings.IndexByte(base, enc[i]) * mult)
	}
	return nb
}
