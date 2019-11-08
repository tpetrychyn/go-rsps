package writer

import (
	"bytes"
	"log"
)

type InterfaceTextPacketWriter struct {
	interfaceId int
	text string
}

func (i *InterfaceTextPacketWriter) Build() []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte{126})
	buf.Write([]byte{0, byte(len(i.text) + 3)})
	buf.Write([]byte(i.text))
	buf.Write([]byte{10})
	buf.Write([]byte{byte(i.interfaceId << 8), byte(i.interfaceId)})
	log.Printf("Writing interface: %+v", buf.Bytes())
	return buf.Bytes()
}

