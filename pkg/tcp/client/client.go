package client

import (
	"net"
)

type TcpClient struct {
	conn          net.Conn
	rate          int
	rateUpdatesCh chan int
	counter       int
}

// Close cleans up the information that TcpClient holds.
// The rateUpdatesCh channel is closed later than the tcp connection to prevent
// sending messages to a closed channel during the listening of rate updates, which would cause panic.
func (tcpClient *TcpClient) Close() {
	tcpClient.conn.Close()
	close(tcpClient.rateUpdatesCh)
}

// NewTcpClient creates a new TcpClient with a tcp connection and an initialCounter,
// which specifies the starting number to send to the server via the connection.
func NewTcpClient(conn net.Conn, initialCounter int) *TcpClient {
	return &TcpClient{
		conn:          conn,
		rate:          0, // By default, no messages will be fired from client to server
		rateUpdatesCh: make(chan int),
		counter:       initialCounter,
	}
}
