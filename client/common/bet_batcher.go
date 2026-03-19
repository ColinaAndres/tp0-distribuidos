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
	batchSizer      BatchSizer
	maxBatchAmount  int
	pendingBet      *Bet
}

type BatchSizer func(currentCount int, bet *Bet) int

func NewBetBatcher(agency string, reader *bufio.Scanner, betSizer BatchSizer, maxBatchAmount int) *BetBatcher {
	return &BetBatcher{
		agency:          agency,
		reader:          reader,
		actualBatch:     make([]Bet, 0),
		actualBatchSize: 0,
		batchSizer:      betSizer,
		maxBatchAmount:  maxBatchAmount,
		pendingBet:      nil,
	}
}

func (betBatcher *BetBatcher) GetBatch() []Bet {
	for len(betBatcher.actualBatch) < betBatcher.maxBatchAmount && betBatcher.reader.Scan() {
		line := betBatcher.reader.Text()
		bet := NewBetFromLine(betBatcher.agency, line)
		batchSize := betBatcher.batchSizer(betBatcher.actualBatchSize, bet)
		if batchSize > maxBatchSize {
			betBatcher.pendingBet = bet
			break
		}
		betBatcher.actualBatch = append(betBatcher.actualBatch, *bet)
		betBatcher.actualBatchSize = batchSize
	}
	result := betBatcher.actualBatch
	betBatcher.restart()
	return result
}

func (betBatcher *BetBatcher) restart() {
	betBatcher.actualBatch = make([]Bet, 0)
	betBatcher.actualBatchSize = 0
	if betBatcher.pendingBet != nil {
		betBatcher.actualBatch = append(betBatcher.actualBatch, *betBatcher.pendingBet)
		betBatcher.actualBatchSize = betBatcher.batchSizer(betBatcher.actualBatchSize, betBatcher.pendingBet)
		betBatcher.pendingBet = nil
	}
}
