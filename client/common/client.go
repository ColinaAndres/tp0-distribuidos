package common

import (
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
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
	bet := NewBetFromEnv()
	betProtocol, err := NewBetProtocol(c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err)
		return
	}
	c.betProtocol = betProtocol
	defer c.GracefulShutdown()

	if err := c.betProtocol.SendBet(bet); err != nil {
		log.Errorf(
			"action: send_bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	if err := c.betProtocol.ReceiveConfirmation(); err != nil {
		log.Errorf(
			"action: receive_confirmation | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	log.Infof(
		"action: apuesta_enviada | result: success | dni: %v | numero: %v",
		bet.document,
		bet.number,
	)
}
