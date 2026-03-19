package common

import (
	"math"
	"net"
	"time"
)

// Socket is a wrapper around net.Conn that provides methods for sending and receiving data,
// ensuring that all data is sent or received as needed.
type Socket struct {
	conn net.Conn
}

// NewSocket creates a new Socket instance by connecting to the specified server address.
// It returns a pointer to the Socket and any error encountered during the connection process.
func NewSocket(serverAddress string) (*Socket, error) {
	maxAttemps := 3
	var err error
	for i := 0; i < maxAttemps; i++ {
		conn, err_aux := net.Dial("tcp", serverAddress)
		if err_aux == nil {
			return &Socket{conn}, nil
		}
		err = err_aux

		// Exponential backoff strategy for reconnection attempts
		waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
		time.Sleep(waitTime)
	}
	return nil, err
}

// SendAll sends the entire byte slice data to the connected server.
// It ensures that all bytes are sent.
// If an error occurs during sending, it returns the error.
func (s *Socket) SendAll(data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		sended, err := s.conn.Write(data[totalSent:])
		if err != nil {
			return err
		}
		totalSent += sended
	}
	return nil
}

// ReceiveAll receives a specified amount of bytes from the connected server.
// It ensures that the exact amount of bytes is received and returns the data as a byte slice.
// If an error occurs during receiving, it returns the error.
func (s *Socket) ReceiveAll(amountToReceive int) ([]byte, error) {
	buffer := make([]byte, amountToReceive)
	total_received := 0
	for total_received < amountToReceive {
		received, err := s.conn.Read(buffer[total_received:])
		if err != nil {
			return nil, err
		}
		total_received += received
	}
	return buffer, nil
}

// Close closes the connection to the server.
// It returns any error encountered during the closing process.
func (s *Socket) Close() error {
	return s.conn.Close()
}
