package net

import (
	"context"
	"github.com/gtank/isaac"
	"log"
	"net"
	"rsps/entity"
	"rsps/net/packet"
)

const LoginState = "LOGIN_STATE"
const ConnectionStatus = "CONNECTION_STATUS"
const Disconnected = "DISCONNECTED"

type Connection struct {
	TCPConn     *net.TCPConn
	Context     context.Context
	PacketQueue []*packet.Packet
	Player      *entity.Player
	Encryptor   *isaac.ISAAC
	Decryptor   *isaac.ISAAC
}

func (c *Connection) SetValue(key interface{}, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}

func (c *Connection) GetValue(key interface{}) interface{} {
	return c.Context.Value(key)
}

func (c *Connection) WriteFrame(id int, b []int) {
	buf := make([]byte, len(b))
	for k, v := range b {
		if v < 0 {
			buf[k] = byte(256 + v)
		} else {
			buf[k] = byte(v)
		}
	}

	log.Printf("Writing Packet id: %d, %+v, %s", id, buf, string(buf))
	buf = append([]byte{byte(id)}, buf...)
	c.TCPConn.Write(buf)
}

func (c *Connection) Wb(b []byte) {
	//log.Printf("Writing Packet: %+v", b)
	c.TCPConn.Write(b)
}

func (c *Connection) W(b []int) {
	buf := make([]byte, len(b))
	for k, v := range b {
		if v < 0 {
			buf[k] = byte(256 + v)
		} else {
			buf[k] = byte(v)
		}
	}

	log.Printf("Writing Packet: %+v", buf)
	c.TCPConn.Write(buf)
}
