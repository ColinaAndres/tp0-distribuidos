package common

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

const (
	batchDivider     = "|"
	betDivider       = ","
	confirmation     = 1
	lengthPrefixSize = 2
	confirmationSize = 1
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
	serializedBet := serializeBet(bet)
	return betProtocol.skt.SendAll(serializedBet)
}

// ReceiveConfirmation waits for a confirmation byte from the server after sending a bet.
// It returns an error if the confirmation is not received or if any error occurs during receiving.
func (betProtocol *BetProtocol) ReceiveConfirmation() error {
	buff, err := betProtocol.skt.ReceiveAll(confirmationSize)
	if err != nil {
		return err
	} else if buff[0] != confirmation {
		return fmt.Errorf("confirmation not received")
	}
	return nil
}

// LookAheadBatchSize calculates the size of a batch of bets if a new bet is added to it.
func (betProtocol *BetProtocol) LookAheadBatchSize(currentSize int, bet *Bet) int {
	if currentSize == 0 {
		return lengthPrefixSize + len(serializeBet(bet))
	}
	return currentSize + len(batchDivider) + len(serializeBetPayload(bet))
}

// Close closes the Socket connection used by the BetProtocol.
func (betProtocol *BetProtocol) Close() error {
	return betProtocol.skt.Close()
}

// Aux function to serialize a Bet struct into a byte slice using CSV format.
func serializeBet(bet *Bet) []byte {
	serializedBet := serializeBetPayload(bet)
	return append(serializeLengthPrefix(serializedBet), serializedBet...)
}

// Aux function to serialize a Bet struct into a byte slice using CSV format.
func serializeBetPayload(bet *Bet) []byte {
	csvBet := strings.Join([]string{bet.agency, bet.name, bet.lastName, bet.document, bet.birth, bet.number}, betDivider)
	return []byte(csvBet)
}

// Aux function to serialize the length of the data as a 2-byte big-endian prefix.
func serializeLengthPrefix(data []byte) []byte {
	length := len(data)
	if length > math.MaxUint16 {
		panic("data length exceeds maximum allowed size")
	}
	prefix := make([]byte, lengthPrefixSize)
	binary.BigEndian.PutUint16(prefix, uint16(length))
	return prefix
}
