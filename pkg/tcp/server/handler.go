package server

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type MessageRateSet struct {
	Rate int `json:"rate"`
}

type Handler struct {
	conn      net.Conn
	done      chan struct{}
	closeOnce sync.Once

	varints chan int64
}

// NewHandler creates a new handler for the tcp connection to receive counters
// and broadcast rate updates later on.
func NewHandler(conn net.Conn) *Handler {
	return &Handler{
		conn:    conn,
		varints: make(chan int64),
		done:    make(chan struct{}),
	}
}

// Serve listens to the counters fired from client and validates them.
// Serve also periodically (every 5 seconds) sends the new rate [0,100) to clients.
func (handler *Handler) Serve() {
	defer handler.close()
	go handler.read()

	tickerSet, tickerValidate := time.NewTicker(time.Second*5), time.NewTicker(time.Second)
	rate, isRateChanged := 0, true
	lastVarint, varintsReceivedCount := int64(0), 0

	for {
		select {
		case <-tickerValidate.C:
			if isRateChanged {
				setVarintsReceivedCount(&varintsReceivedCount, 0)
				setIsRateChanged(&isRateChanged, false)
				continue
			}

			if varintsReceivedCount != rate {
				log.Printf("Not equal: rate=%v varintsReceivedCount=%v", rate, varintsReceivedCount)
				return
			}

			setVarintsReceivedCount(&varintsReceivedCount, 0)

		case <-tickerSet.C:
			rate = rand.Intn(100)

			log.Println("Set rate:", rate)
			err := handler.write(rate)

			if err != nil {
				log.Println("Error broadcasting new rate update to client. Message:", err)
				return
			}

			setVarintsReceivedCount(&varintsReceivedCount, 0)
			setIsRateChanged(&isRateChanged, true)

		case nextVarint := <-handler.varints:
			// Whenever we receive varints, check if varint is in correct sequence.
			// In other words, client must write sequentially to server.
			if nextVarint != lastVarint+1 {
				log.Printf(
					"Out of order: lastVarint=%v nextVarint=%v",
					lastVarint, nextVarint,
				)

				return
			}

			lastVarint = nextVarint
			varintsReceivedCount++

		case <-handler.done:
			return
		}
	}
}

// close cleans up handler, such as closing the connection and channels.
func (handler *Handler) close() {
	handler.closeOnce.Do(func() {
		close(handler.done)

		err := handler.conn.Close()
		if err != nil {
			log.Println("Error closing connection. Message:", err)
		}

		log.Println("Disconnected", handler.conn.RemoteAddr())
	})
}

// read listens to the counters fired from clients and passes it for validation.
func (handler *Handler) read() {
	defer handler.close()

	reader := bufio.NewReader(handler.conn)

	for {
		varint, err := binary.ReadVarint(reader)
		if err != nil {
			log.Println("Error reading varint from client. Message:", err)
			return
		}

		handler.varints <- varint
	}
}

// write fires the new rate update in JSON format to the client.
func (handler *Handler) write(rate int) error {
	message := MessageRateSet{
		Rate: rate,
	}

	encoded, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = handler.conn.Write(encoded)
	if err != nil {
		return err
	}

	log.Printf("Broadcasted new rate of %d to %s.", rate, handler.conn.RemoteAddr())
	return nil
}

// setVarintsReceivedCount will reset the the number of varints we have received every second and
// prevent varintsReceivedCount from being an evergrowing number, so that varintsReceivedCount
// will solely be the number of varints we have received every second.
func setVarintsReceivedCount(varintsReceivedCount *int, newVarintsReceivedCount int) {
	*varintsReceivedCount = newVarintsReceivedCount
}

// isRateChanged is our flag to determine if we had a change of rate, and this
// flag will be used to skip validations depending on the use case.
func setIsRateChanged(isRateChanged *bool, newIsRateChanged bool) {
	*isRateChanged = newIsRateChanged
}
