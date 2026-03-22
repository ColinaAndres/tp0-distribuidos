package common

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

const (
	batchDivider       = "|"
	betDivider         = ","
	documentDivider    = ","
	finalizationByte   = 0
	batchSendingByte   = 1
	winnersRequestByte = 2
	confirmation       = 1
	lengthPrefixSize   = 2
	confirmationSize   = 1
	batchHeaderSize    = 1 + lengthPrefixSize
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

// SendBatch sends a batch of Bet. It first serializes the batch
// into a byte slice, including a length prefix and batch dividers, and then sends it
// It returns an error if any issue occurs during the sending process.
func (betProtocol *BetProtocol) SendBatch(bets []Bet) error {
	var serializedBatch []byte
	for i, bet := range bets {
		if i > 0 {
			serializedBatch = append(serializedBatch, []byte(batchDivider)...)
		}
		serializedBatch = append(serializedBatch, serializeBetPayload(&bet)...)
	}
	serializedBatch = append(serializeBatchHeader(serializedBatch), serializedBatch...)
	return betProtocol.skt.SendAll(serializedBatch)
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
		return batchHeaderSize + len(serializeBetPayload(bet))
	}
	return currentSize + len(batchDivider) + len(serializeBetPayload(bet))
}

// SendFinalization sends a finalization byte to the server to
// indicate that no more bets will be sent.
func (betProtocol *BetProtocol) SendFinalization() error {
	return betProtocol.skt.SendAll([]byte{finalizationByte})
}

func (betProtocol *BetProtocol) ReceiveWinners() ([]string, error) {
	if err := betProtocol.skt.SendAll([]byte{winnersRequestByte}); err != nil {
		return nil, err
	}
	lengthPrefixBytes, err := betProtocol.skt.ReceiveAll(lengthPrefixSize)
	if err != nil {
		return nil, err
	}
	winnersLength := int(binary.BigEndian.Uint16(lengthPrefixBytes))
	winnersData, err := betProtocol.skt.ReceiveAll(winnersLength)
	if err != nil {
		return nil, err
	}
	winners := strings.Split(string(winnersData), documentDivider)
	return winners, nil
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

func serializeBatchHeader(batch []byte) []byte {
	return append([]byte{batchSendingByte}, serializeLengthPrefix(batch)...)
}
