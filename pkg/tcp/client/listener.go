package client

import (
	"encoding/json"
	"log"
)

type RateUpdateBroadcast struct {
	Rate int `json:"rate"`
}

// HandleNewRateUpdates listens and reads the JSON rate update broadcast from the server.
// It will block and wait until we get a response from the server.
// Upon getting the new rate update from server, the new rate will be communicated via
// rateUpdatesCh channel to our main thread to send packets with the new rate
func (tcpClient *TcpClient) HandleNewRateUpdates() {
	// Instead of using json Decoder, we will split the reading and unmarshalling process.
	// This is so that we can have clearer logs and pinpoint the point of failures more easily,
	// making the debugging process easier in the future.
	buffer := make([]byte, 1024)

	for {
		n, err := tcpClient.conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from server. Message:", err)
			return
		}

		var rateUpdate RateUpdateBroadcast
		if err := json.Unmarshal(buffer[:n], &rateUpdate); err != nil {
			log.Println("Error decoding JSON rate update broadcast from server. Message:", err)
			return
		}

		newRate := rateUpdate.Rate
		tcpClient.rateUpdatesCh <- newRate
		log.Printf("Got a new rate update of %d!", newRate)
	}
}
