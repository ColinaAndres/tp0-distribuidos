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
		n, err := s.conn.Write(data[totalSent:])
		if err != nil {
			return err
		}
		totalSent += n
	}
	return nil
}
