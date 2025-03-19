package common

import (
	"time"
	"os"
	"os/signal"
	"syscall"

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
	config 	ClientConfig
	sigChan chan os.Signal
	central CentralLoteriaNacional
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, central CentralLoteriaNacional) *Client {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	client := &Client{
		config: config,
		sigChan: sigChan,
		central: central,
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) SendBet(bet Bet) {
	select {
		case <- c.sigChan:
			log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
		default:
			err := c.central.SendBet(bet)
			if err != nil {
				log.Criticalf(
					"action: connect | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
			}
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet.Document, bet.Number)
	}
}