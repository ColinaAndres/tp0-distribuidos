package common

type BetProtocol struct {
	skt *Socket
}

func NewBetProtocol(serverAddress string) (*BetProtocol, error) {
	skt, err := NewSocket(serverAddress)
	if err != nil {
		return nil, err
	}
	return &BetProtocol{skt: skt}, nil
}
