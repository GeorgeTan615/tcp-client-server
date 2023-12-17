package client

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConnectionMock struct {
	net.Conn
}

func (m *ConnectionMock) Close() error {
	return nil
}

func newConnectionMock() *ConnectionMock {
	return &ConnectionMock{}
}

func TestNewTcpClient(t *testing.T) {
	mockConn := newConnectionMock()
	initialCounter := 1

	tcpClient := NewTcpClient(mockConn, initialCounter)

	assert.Equal(t, mockConn, tcpClient.conn)
	assert.Equal(t, initialCounter, tcpClient.counter)
	assert.Equal(t, 0, tcpClient.rate)
	assert.NotEqual(t, nil, tcpClient.rateUpdatesCh)
}

func TestClose(t *testing.T) {
	mockConn := newConnectionMock()
	initialCounter := 1

	tcpClient := NewTcpClient(mockConn, initialCounter)

	tcpClient.Close()

	// Check if rateUpdatesCh is closed
	select {
	case _, ok := <-tcpClient.rateUpdatesCh:
		assert.False(t, ok)
	default:
	}
}
