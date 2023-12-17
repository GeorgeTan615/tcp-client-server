package client

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

// SendCounters will send the varint packets to the server
// every second according to the current rate.
// SendCounters will also actively listen to rate updates
// from the server and adjust the firing rate accordingly.
func (tcpClient *TcpClient) SendCounters() {
	go tcpClient.HandleNewRateUpdates()

	tickerSend := time.NewTicker(time.Second)
	defer tickerSend.Stop()

	for {
		select {
		case newRate, ok := <-tcpClient.rateUpdatesCh:
			if !ok {
				log.Println("Error occured reading from rateUpdatesCh")
				return
			}

			tcpClient.rate = newRate
			log.Println("Updated new rate:", tcpClient.rate)

		case <-tickerSend.C:
			log.Printf("Sending %d packets. Current counter: %d", tcpClient.rate, tcpClient.counter)
			err := sendVarintPackets(tcpClient)

			if err != nil {
				log.Println("Error sending packets. Message:", err)
				return
			}
		}
	}
}

// sendVarintPackets will fire the varints to the server.
func sendVarintPackets(tcpClient *TcpClient) error {
	// To minimize our network I/O cost,
	// we will consolidate all the varints and only make one single write,
	// this prevents multiple writes over the tcp connection for each varint.
	varintBuffer := new(bytes.Buffer)

	for i := 0; i < tcpClient.rate; i++ {
		varintTempBuffer := make([]byte, binary.MaxVarintLen64)
		bytesWritten := binary.PutVarint(varintTempBuffer, int64(tcpClient.counter))
		varintBuffer.Write(varintTempBuffer[:bytesWritten])
		tcpClient.counter++
	}

	_, err := tcpClient.conn.Write(varintBuffer.Bytes())
	return err
}
