package common

import "bufio"

const (
	maxBatchSize = 8192
)

type BetBatcher struct {
	agency          string
	reader          *bufio.Scanner
	actualBatch     []Bet
	actualBatchSize int
	betSizer        BetSizer
}

type BetSizer func(currentCount int, bet Bet) int

func NewBetBatcher(agency string, reader *bufio.Scanner, betSizer BetSizer) *BetBatcher {
	return &BetBatcher{
		agency:          agency,
		reader:          reader,
		actualBatch:     make([]Bet, 0),
		actualBatchSize: 0,
		betSizer:        betSizer,
	}
}
