package packet

import "math"

type Stream struct {
	bitBuffer map[uint]byte
	bitPosition uint
}

func NewStream() *Stream {
	bitBuffer := make(map[uint]byte)
	bitBuffer[0] = 0
	return &Stream{
		bitBuffer:   bitBuffer,
		bitPosition: 0,
	}
}

func (s *Stream) Flush() []byte {
	buf := make([]byte, 0)
	for _, v := range s.bitBuffer {
		buf = append(buf, v)
	}
	// TODO: stupid iteration order hack to append to buffer in order!
	for k, v := range s.bitBuffer {
		buf[k] = v
	}

	s.bitBuffer = make(map[uint]byte)
	s.bitBuffer[0] = 0
	s.bitPosition = 0

	return buf
}

func (s *Stream) WriteBits(numBits uint, value uint) {
	if numBits == 0 {
		return
	}

	if numBits > 8 {
		s.WriteBits(numBits - 8, value >> 8)
		numBits = 8
		value >>= numBits - 8
	}
	bytePos := uint(len(s.bitBuffer) - 1)
	if numBits+s.bitPosition > 8 {
		diff := uint(math.Abs(float64(int(numBits - s.bitPosition))))

		shift := numBits-(8-s.bitPosition)
		s.bitBuffer[bytePos] = byte(uint(s.bitBuffer[bytePos])<<(8-s.bitPosition) | value>>shift)
		bytePos++
		shift = 8-shift
		s.bitBuffer[bytePos] = byte(value << shift)
		s.bitPosition = 8 - diff
	} else {
		s.bitBuffer[bytePos] = byte(uint(s.bitBuffer[bytePos])<<numBits | value >> s.bitPosition)
		s.bitPosition += numBits
	}
}
