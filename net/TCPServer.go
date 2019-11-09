package net

import (
	"log"
	"net"
	"strconv"
)

type TcpServer struct {
	Port         int
	Listener     net.Listener
}

func NewTcpServer(port int) *TcpServer {
	return &TcpServer{
		Port:         port,
	}
}

func (server *TcpServer) Start() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	log.Printf("Local channel bound at %v \n", server.Port)

	l := &UpstreamLoginHandler{}

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		client := NewTcpClient(connection, l)

		go client.Read()
		go client.Write()
		go client.ProcessUpstream()
		go client.Tick()
	}
}

func (server *TcpServer) Stop() {
	if server.Listener != nil {
		_ = server.Listener.Close()
	}
}
