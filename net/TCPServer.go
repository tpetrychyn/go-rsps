package net

import (
	"log"
	"net"
	"rsps/entity"
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

var clients = make([]*TCPClient, 0)

func (server *TcpServer) Start() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	log.Printf("Local channel bound at %v \n", server.Port)
	world := entity.WorldProvider()
	go world.Tick()

	l := &LoginHandler{}

	//go func() {
	//	for {
	//		<-time.After(600 * time.Millisecond)
	//		for _, c := range clients {
	//			c.Tick()
	//		}
	//		for _, c := range clients {
	//			c.Write()
	//		}
	//	}
	//}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		client := NewTcpClient(connection, l, world)

		go client.Read()
		go client.Write()
		go client.ProcessUpstream()
		clients = append(clients, client)
		go client.Tick()
	}
}

func (server *TcpServer) Stop() {
	if server.Listener != nil {
		_ = server.Listener.Close()
	}
}
