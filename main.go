package main

import (
	"rsps/net"
	"rsps/util"
)

func main() {
	util.LoadItemDefinitions()

	server := net.NewTcpServer(43594)
	server.Start()
}
