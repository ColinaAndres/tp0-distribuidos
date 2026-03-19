package common

import "net"

type Socket struct {
	conn net.Conn
}

func NewSocket(serverAddress string) (*Socket, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	v := &Socket{conn}

	return v, nil

}

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

func (s *Socket) Close() error {
	return s.conn.Close()
}
