package common

import (
	"bufio"
	"os"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	MaxBatchAmount int
	DataRoute      string
}

// Client Entity that encapsulates how
type Client struct {
	config      ClientConfig
	betProtocol *BetProtocol
	running     bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// GracefulShutdown Closes the client socket if it is open.
func (c *Client) GracefulShutdown() {
	c.running = false
	if c.betProtocol != nil {
		log.Infof("action: closing_client_socket | result: in_progress | client_id: %v", c.config.ID)
		err := c.betProtocol.Close()
		if err != nil {
			log.Errorf("action: closing_client_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
		} else {
			log.Infof("action: closing_client_socket | result: success | client_id: %v", c.config.ID)
		}
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClient() {
	c.running = true

	if !c.initializeProtocol() {
		return
	}
	defer c.GracefulShutdown()

	file, ok := c.openDataFile()
	if !ok {
		return
	}
	defer file.Close()

	c.processBatches(file)
	c.informWinners()
}

// initializeProtocol Initializes the protocol used to send bets to the server
func (c *Client) initializeProtocol() bool {
	betProtocol, err := NewBetProtocol(c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return false
	}

	c.betProtocol = betProtocol
	return true
}

// openDataFile Opens the file that contains the bets to be sent to the server
func (c *Client) openDataFile() (*os.File, bool) {
	file, err := os.Open(c.config.DataRoute)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return nil, false
	}

	return file, true
}

// processBatches Reads the bets from the file and sends them to the server in batches
func (c *Client) processBatches(file *os.File) {
	reader := bufio.NewScanner(file)
	betBatcher := NewBetBatcher(
		c.config.ID,
		reader,
		c.betProtocol.LookAheadBatchSize,
		c.config.MaxBatchAmount,
	)

	for betsToSend := betBatcher.GetBatch(); betsToSend != nil; betsToSend = betBatcher.GetBatch() {
		if !c.sendBets(betsToSend) {
			return
		}
		if !c.waitConfirmation(betsToSend) {
			return
		}

	}
	c.SendFinalization()
}

// sendBets Sends a batch of bets to the server using the protocol
func (c *Client) sendBets(bets []Bet) bool {
	if err := c.betProtocol.SendBatch(bets); err != nil {
		log.Errorf(
			"action: send_batch | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return false
	}
	return true
}

// waitConfirmation Waits for the server to confirm that the batch of bets was received and processed
func (c *Client) waitConfirmation(bets []Bet) bool {
	if err := c.betProtocol.ReceiveConfirmation(); err != nil {
		log.Errorf(
			"action: receive_confirmation | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return false
	}

	log.Infof(
		"action: batch_sending | result: success | client_id: %v | bets_sent: %v",
		c.config.ID,
		len(bets),
	)
	return true
}

func (c *Client) SendFinalization() bool {
	if err := c.betProtocol.SendFinalization(); err != nil {
		log.Errorf(
			"action: send_finalization | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return false
	}
	return true
}

func (c *Client) informWinners() {
	winners, err := c.betProtocol.ReceiveWinners()
	if err != nil {
		log.Errorf(
			"action: receive_winners | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	log.Infof(
		"action: consulta_ganadores | result: success | cant_ganadores: %v",
		len(winners),
	)
}
