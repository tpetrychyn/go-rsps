package net

import (
	"log"
	"net"
	"rsps/entity"
	"rsps/repository"
	"strconv"
	"sync"
	"time"
)

type TcpServer struct {
	Port    int
	Clients *sync.Map
	Listener net.Listener
}

func NewTcpServer(port int) *TcpServer {
	return &TcpServer{
		Port:    port,
		Clients: new(sync.Map),
	}
}

func (server *TcpServer) Start(playerRepository *repository.PlayerRepository) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	log.Printf("Local channel bound at %v \n", server.Port)
	world := entity.WorldProvider()

	l := &LoginHandler{}

	tickGroup := new(sync.WaitGroup)
	updateGroup := new(sync.WaitGroup)

	go func() {
		for {
			// TODO: Could implement a tickTime at the start of each loop and fire when it hits 600ms
			// so that the ticks have the entire 600ms to process instead of rushing at the end
			<-time.After(600 * time.Millisecond)
			// let all clients tick in parallel threads (handle movement and pickup, etc)
			// parallel threads minimizes advantage of pID
			//tickGroup.Add(len(server.Clients))
			world.Tick()
			server.Clients.Range(func(key, value interface{}) bool {
				client := value.(*TCPClient)
				if client.loginState == IngameStage {
					tickGroup.Add(1)
					go client.Tick(tickGroup)
				}
				if client.loginState == Disconnected {
					client.connection.Close()
					server.Clients.Delete(key)
				}
				return true
			})
			tickGroup.Wait()

			// after all have ticked, issue the update packets in parallel
			server.Clients.Range(func(key, value interface{}) bool {
				client := value.(*TCPClient)
				if client.loginState == IngameStage {
					updateGroup.Add(1)
					go client.UpdatePacket(updateGroup)
					go client.PlayerRepository.Save(client.Player.Name, client.Player.Position)
				}
				return true
			})
			updateGroup.Wait()

			server.Clients.Range(func(key, value interface{}) bool {
				client := value.(*TCPClient)
				if client.loginState == IngameStage {
					client.Player.PostUpdate()
				}
				return true
			})

			world.PostUpdate()
		}
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		client := NewTcpClient(connection, l, playerRepository, world)

		go client.Read()
		go client.Write()
		go client.ProcessUpstream()
		server.Clients.Store(client.connection.RemoteAddr().String(), client)
	}
}

func (server *TcpServer) Stop() {
	if server.Listener != nil {
		_ = server.Listener.Close()
	}
}
