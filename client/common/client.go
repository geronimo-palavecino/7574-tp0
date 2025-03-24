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
func (c *Client) Lottery() {
	// The connection is opened
	err := c.central.CreateSocket()
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	
	allBetsSent := false
	for true {
		select {
			case <- c.sigChan:
				c.repo.Close()
				c.central.Close()
				log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
				return
			default:
				if allBetsSent {
					c.repo.Close()
					getResults(c.config.ID, c.central)
					c.central.Close()
					return
				} else {
					var err error
					allBetsSent, err = sendBets(c.repo, c.config.Batch, c.config.ID, c.central)
					if err != nil {
						c.repo.Close()
						c.central.Close()
						return
					}
				}
		}
	}
}

func sendBets(repo *BetRepository, amount int, id int, central CentralLoteriaNacional) (bool, error) {
	bets, err := repo.GetBets(amount, id)
	if err != nil {
		log.Criticalf(
			"action: read_bets | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
		return false, err
	}

	if len(bets) == 0 {
		return true, nil
	}

	err = central.SendBets(bets)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
		return false, err
	}

	log.Infof("action: apuesta_enviada | result: success | cantidad: %v", len(bets))

	return false, nil
}

func getResults(id int, central CentralLoteriaNacional) {
	winners, err := central.GetWinners(id)
	if err != nil {
		log.Criticalf(
			"action: consulta_ganadores| result: fail | client_id: %v | error: %v",
			id,
			err,
		)
	} else {
		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(winners))
	}
}