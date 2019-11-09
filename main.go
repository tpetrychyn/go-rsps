package main

import (
	"log"
	"net"
	"rsps/entity"
	net2 "rsps/net"

	//"rsps/net"
)

const PORT = "43594"

func main() {

	connectionHandler := net2.NewConnectionHandler()

	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":43594")
	ln, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Listening on %s", PORT)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			// TODO: Handle error
		}
		player := entity.NewPlayer()
		go player.Tick()
		go connectionHandler.Listener(player.TCPConn)
		go connectionHandler.writer(newConn)
	}

	//
	//connectionHandler.Listen()

}
