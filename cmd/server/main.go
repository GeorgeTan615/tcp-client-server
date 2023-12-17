package main

import (
	"log"
	"net"
	"test-task/pkg/tcp/server"
)

const (
	address = "localhost:54321"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error creating TCP listener for %s. Message: %v", address, err)
	}

	log.Println("Waiting for new connections on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error accepting new connection. Message: ", err)
		}

		log.Printf("Connected with %s!", conn.RemoteAddr())

		go server.NewHandler(conn).Serve()
	}
}
