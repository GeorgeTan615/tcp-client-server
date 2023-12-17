package client

import (
	"encoding/json"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ConnectionMockWithRateUpdate struct {
	net.Conn
	rateUpdate *RateUpdateBroadcast
}

func newConnectionMockWithRateUpdate(rateUpdate *RateUpdateBroadcast) *ConnectionMockWithRateUpdate {
	return &ConnectionMockWithRateUpdate{
		rateUpdate: rateUpdate,
	}
}

func (m *ConnectionMockWithRateUpdate) Read(b []byte) (int, error) {
	// Simulate the convertion of JSON from server
	bytes, err := json.Marshal(m.rateUpdate)

	if err != nil {
		log.Fatalln(err)
	}

	n := len(bytes)
	copy(b, bytes)
	return n, nil
}

func TestHandleNewRateUpdates(t *testing.T) {
	rateUpdate := &RateUpdateBroadcast{Rate: 10}
	mockConn := newConnectionMockWithRateUpdate(rateUpdate)
	tcpClient := NewTcpClient(mockConn, 1)
	go tcpClient.HandleNewRateUpdates()

	select {
	case receivedRate := <-tcpClient.rateUpdatesCh:
		assert.Equal(t, rateUpdate.Rate, receivedRate)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for rate update")
	}
}
