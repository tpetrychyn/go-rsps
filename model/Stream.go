package model

type Stream struct {
	buffer        []byte
	bitPosition   uint
	currentOffset uint
}

var bitmask = [32]uint{}

func NewStream() *Stream {
	buffer := append(make([]byte, 0), 0)

	for i := 0; i < 32; i++ {
		bitmask[i] = (1 << uint(i)) - 1
	}

	return &Stream{
		buffer:        buffer,
		bitPosition:   0,
		currentOffset: 0,
	}
}

func (s *Stream) Flush() []byte {
	buf := s.buffer[:]

	s.buffer = append(make([]byte, 0), 0)
	s.bitPosition = 0
	s.currentOffset = 0

	return buf
}

func (s *Stream) Write(bytes []byte) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, bytes...)
	} else {
		s.buffer = bytes
	}

	s.currentOffset += uint(len(bytes))
}

func (s *Stream) WriteWordLE(value uint) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = byte(value)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value >> 8)
	s.currentOffset++
}

func (s *Stream) WriteWord(value uint) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = byte(value >> 8)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value)
	s.currentOffset++
}

func (s *Stream) WriteWordA(value uint) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = byte(value >> 8)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value + 128)
	s.currentOffset++
}

func (s *Stream) WriteWordBEA(value uint) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = byte(value + 128)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value >> 8)
	s.currentOffset++
}

func (s *Stream) WriteByte(value byte) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = value
	s.currentOffset++
}

func (s *Stream) WriteInt(value int) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = byte(value >> 24)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value >> 16)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value >> 8)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value)
	s.currentOffset++
}

func (s *Stream) WriteDWord_v1(value int) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = byte(value >> 8)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value >> 24)
	s.currentOffset++
	s.buffer = append(s.buffer, 0 )
	s.buffer[s.currentOffset] = byte(value >> 16)
	s.currentOffset++
}

func (s *Stream) WriteDWord_v2(value int) {
	if s.currentOffset > 0 {
		s.buffer = append(s.buffer, 0)
	}

	s.buffer[s.currentOffset] = byte(value >> 16)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value >> 24)
	s.currentOffset++
	s.buffer = append(s.buffer, 0)
	s.buffer[s.currentOffset] = byte(value)
	s.currentOffset++
	s.buffer = append(s.buffer, 0 )
	s.buffer[s.currentOffset] = byte(value >> 8)
	s.currentOffset++
}

func (s *Stream) SkipByte() {
	s.currentOffset++
}

func (s *Stream) WriteBits(numBits uint, value uint) {
	if numBits == 0 {
		return
	}

	if numBits > 8 {
		s.WriteBits(numBits-8, value>>8)
		numBits = 8
		value >>= numBits - 8
	}
	bytePos := uint(len(s.buffer) - 1)
	if numBits+s.bitPosition > 8 {
		shift := numBits - (8 - s.bitPosition)
		s.buffer[bytePos] = byte(uint(s.buffer[bytePos]) | (value>>shift)&bitmask[8-int(shift)])
		bytePos++
		s.buffer = append(s.buffer, 0)
		shift = 8 - shift
		s.buffer[bytePos] = byte(value << shift)
		s.bitPosition = 8 - shift
	} else {
		s.buffer[bytePos] = byte(uint(s.buffer[bytePos]) | (value << (8 - s.bitPosition - numBits)))
		s.bitPosition += numBits
	}
}
