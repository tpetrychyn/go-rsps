package packet

import (
	"testing"
)

func TestWriteBits(t *testing.T) {
	s := NewStream()

	s.WriteBits(7, 14)
	s.WriteBits(3, 7)
	if string([]byte{29, 192}) != string(s.Flush()) {
		t.Fatalf("Expected: %+v, got: %+v", []byte{29, 192}, s.Flush())
	}

	s.WriteBits(7, 14)
	s.WriteBits(3, 7)
	if string([]byte{29, 192}) != string(s.Flush()) {
		t.Fatalf("Expected: %+v, got: %+v", []byte{29, 192}, s.Flush())
	}
}

func TestStream_WriteBitsLargerThanByte(t *testing.T) {
	s := NewStream()
	s.WriteBits(7, 0)
	s.WriteBits(11, 2047)

	if string([]byte{1, 127, 240}) != string(s.Flush()) {
		t.Fatalf("Expected: %+v, got: %+v", []byte{1, 127, 240}, s.Flush())
	}
}
