package packet

import (
	"testing"
)

func TestWriteBits(t *testing.T) {
	s := NewStream()

	s.WriteBits(7, 14)
	s.WriteBits(3, 7)

	flushed := s.Flush()
	if string([]byte{29, 192}) != string(flushed) {
		t.Fatalf("Expected: %+v, got: %+v", []byte{29, 192}, flushed)
	}
}

func TestStream_WriteBitsLargerThanByte(t *testing.T) {
	s := NewStream()
	s.WriteBits(7, 0)
	s.WriteBits(11, 2047)

	flushed := s.Flush()
	if string([]byte{1, 255, 192}) != string(flushed) {
		t.Fatalf("Expected: %+v, got: %+v", []byte{1, 255, 192}, flushed)
	}
}

func TestStream_WriteBitsSmallerThanOneByte(t *testing.T) {
	s := NewStream()
	s.WriteBits(1, 1)
	if s.bitBuffer[0] != 128 {
		t.Fatalf("Expected: %v, got: %+v", 128, s.bitBuffer[0])
	}

	s.WriteBits(2, 1)
	if s.bitBuffer[0] != 160 {
		t.Fatalf("Expected: %v, got: %+v", 160, s.bitBuffer[0])
	}

	s.WriteBits(3, 1)
	if s.bitBuffer[0] != 164 {
		t.Fatalf("Expected: %v, got: %+v", 164, s.bitBuffer[0])
	}

	s.WriteBits(1, 1)
	if s.bitBuffer[0] != 166 {
		t.Fatalf("Expected: %v, got: %+v", 166, s.bitBuffer[0])
	}

	flushed := s.Flush()
	if string([]byte{166}) != string(flushed) {
		t.Fatalf("Expected: %+v, got: %+v", []byte{166}, flushed)
	}
}
