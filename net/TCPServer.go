package net

import (
	"log"
	"net"
	"rsps/entity"
	"strconv"
	"sync"
	"time"
)

type TcpServer struct {
	Port     int
	Clients map[string]*TCPClient
	Listener net.Listener
}

func NewTcpServer(port int) *TcpServer {
	return &TcpServer{
		Port: port,
		Clients: make(map[string]*TCPClient),
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
	world := entity.WorldProvider()
	go world.Tick()

	l := &LoginHandler{}

	tickGroup := new(sync.WaitGroup)

	go func() {
		for {
			<-time.After(600 * time.Millisecond)
			// let all clients tick in parallel threads (handle movement and pickup, etc)
			// parallel threads minimizes advantage of pID
			//tickGroup.Add(len(server.Clients))
			for k, c := range server.Clients {
				if c.loginState == IngameStage {
					tickGroup.Add(1)
					go c.Tick(tickGroup)
				}
				if c.loginState == Disconnected {
					server.Clients[k].connection.Close()
					delete(server.Clients, k)
				}
			}
			tickGroup.Wait()

			// after all have ticked, issue the update packets in parallel
			for _, c := range server.Clients {
				if c.loginState == IngameStage {
					go c.UpdatePacket()
				}
			}
		}
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		client := NewTcpClient(connection, l, world)

		go client.Read()
		go client.Write()
		go client.ProcessUpstream()
		server.Clients[client.connection.RemoteAddr().String()] = client
	}
}

func (server *TcpServer) Stop() {
	if server.Listener != nil {
		_ = server.Listener.Close()
	}
}
