package hash

import (
	"testing"
)

func TestEncodeShortenURLID(t *testing.T) {

	urlID := uint64(100)
	// encoding
	want := "1c"
	got := Encode(urlID)

	if want != got {
		t.Errorf("got %s, want %s, given %d", got, want, urlID)
	}
}

func TestDecodeShortURLPath(t *testing.T) {

	short := "1c"
	// decoding
	want := uint64(100)
	got := Decode(short)

	if want != got {
		t.Errorf("got %d, want %d, given %s", got, want, short)
	}
}
