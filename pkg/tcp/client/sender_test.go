package client

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConnectionMockForSend struct {
	net.Conn
	varintBytes []byte
}

func (m *ConnectionMockForSend) Write(b []byte) (n int, err error) {
	m.varintBytes = b // Assign so we can perform comparison later
	return 0, nil
}

func newConnectionMockForSend() *ConnectionMockForSend {
	return &ConnectionMockForSend{}
}

func TestSendVarintPackets(t *testing.T) {
	mockConnection := newConnectionMockForSend()

	initialCounter, rate := 1, 10
	tcpClient := NewTcpClient(mockConnection, initialCounter)
	tcpClient.rate = rate
	sendVarintPackets(tcpClient)

	for i := initialCounter; i < rate+1; i++ {
		value, n := binary.Varint(mockConnection.varintBytes)
		assert.Equal(t, int64(i), value)
		mockConnection.varintBytes = mockConnection.varintBytes[n:]
	}
}
