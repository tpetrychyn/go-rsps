package main

import (
	"rsps/net"
)

func main() {

	connectionHandler := net.NewConnectionHandler()
	connectionHandler.Listen()

}
