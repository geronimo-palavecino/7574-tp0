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
	ID            int
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	Batch 		  int
}

// Client Entity that encapsulates how
type Client struct {
	config 	ClientConfig
	sigChan chan os.Signal
	central CentralLoteriaNacional
	repo 	*BetRepository
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, central CentralLoteriaNacional, repo *BetRepository) *Client {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	client := &Client{
		config: config,
		sigChan: sigChan,
		central: central,
		repo: repo,
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) SendBets() {
	for true {
		select {
			case <- c.sigChan:
				c.repo.Close()
				log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
				return
			default:
				bets, err := c.repo.GetBets(c.config.Batch, c.config.ID)
				if err != nil {
					log.Criticalf(
						"action: read_bets | result: fail | client_id: %v | error: %v",
						c.config.ID,
						err,
					)
					return
				}

				if len(bets) == 0 {
					c.repo.Close()
					return
				}

				err = c.central.SendBets(bets)
				if err != nil {
					log.Criticalf(
						"action: connect | result: fail | client_id: %v | error: %v",
						c.config.ID,
						err,
					)
					return
				}

				log.Infof("action: apuesta_enviada | result: success | cantidad: %v", len(bets))
		}
	}
}