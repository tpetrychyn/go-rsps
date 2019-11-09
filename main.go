package main

import "rsps/net"

func main() {

	//connectionHandler := net.NewConnectionHandler()
	//
	//connectionHandler.Listen()

	server := net.NewTcpServer(43594)

	server.Start()

}
