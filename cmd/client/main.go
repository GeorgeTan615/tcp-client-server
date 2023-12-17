package main

import (
	"log"
	"net"
	"test-task/pkg/tcp/client"
)

const (
	address = "localhost:54321"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Fatalf("Error while attempting to establish TCP connection with %s. Message: %s", err, address)
	}

	log.Println("Connection established with", address)

	initialCounter := 1 // Start at 1 as per the requirements
	tcpClient := client.NewTcpClient(conn, initialCounter)
	defer tcpClient.Close()

	tcpClient.SendCounters()
}
