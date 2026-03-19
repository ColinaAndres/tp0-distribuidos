package common

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	divider      = ","
	confirmation = 1
)

// BetProtocol is a struct that encapsulates the logic for sending bets
// and receiving responses from a server using a custom Socket connection.
type BetProtocol struct {
	skt *Socket
}

// NewBetProtocol creates a new BetProtocol instance with a conection by custom Socket
// to the specified server address. It returns a pointer to the BetProtocol and any
// error encountered during the connection process.
func NewBetProtocol(serverAddress string) (*BetProtocol, error) {
	skt, err := NewSocket(serverAddress)
	if err != nil {
		return nil, err
	}
	return &BetProtocol{skt: skt}, nil
}

// SendBet sends a Bet instance to the server using the BetProtocol's Socket connection.
// It returns an error if any issue occurs during the sending process.
func (betProtocol *BetProtocol) SendBet(bet *Bet) error {
	// TODO: modularizar esta funcion
	serializedBet := []byte(strings.Join([]string{bet.agency, bet.name, bet.lastName, bet.document, bet.birth, bet.number}, divider))
	length_prefix := make([]byte, 2)
	binary.BigEndian.PutUint16(length_prefix, uint16(len(serializedBet)))
	message := append(length_prefix, serializedBet...)

	return betProtocol.skt.SendAll(message)
}

// ReceiveConfirmation waits for a confirmation byte from the server after sending a bet.
// It returns an error if the confirmation is not received or if any error occurs during receiving.
func (betProtocol *BetProtocol) ReceiveConfirmation() error {
	buff, err := betProtocol.skt.ReceiveAll(1)
	if err != nil {
		return err
	} else if buff[0] != confirmation {
		return fmt.Errorf("confirmation not received")
	}
	return nil
}
