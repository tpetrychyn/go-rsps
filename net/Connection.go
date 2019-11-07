package net

import (
	"bytes"
	"context"
	"github.com/gtank/isaac"
	"net"
)

const LoginState = "LOGIN_STATE"

type Connection struct {
	*net.TCPConn
	context.Context
	packetQueue bytes.Buffer
	Encryptor   isaac.ISAAC
	Decryptor   isaac.ISAAC
}

func (c *Connection) SetValue(key interface{}, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}

func (c *Connection) GetValue(key interface{}) interface{} {
	return c.Context.Value(key)
}
