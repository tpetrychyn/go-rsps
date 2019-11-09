package login

import (
	"bufio"
	"encoding/binary"
)

type LoginResponse struct {
	ReturnCode   byte
	PlayerRights byte
	Flagged      byte
}

func (l *LoginResponse) Write(writer *bufio.Writer) {
	binary.Write(writer, binary.BigEndian, l)
}
