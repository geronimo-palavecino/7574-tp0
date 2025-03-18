package common

import (
	"bufio"
	"fmt"
	"net"
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
func NewClient(config ClientConfig) *Client {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	central := NewCentralLoteriaNacional(config.ServerAddress)

	client := &Client{
		config: config,
		sigChan: sigChan,
		central: central,
	}
	return client
}

func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
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
	}
}

// ClientLoop Sen messages to the server until some time threshold is met or a SIGTERM is catched 
/*func (c *Client) ClientLoop() {
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		loopPeriod := time.After(c.config.LoopPeriod)
			select {
				case <- c.sigChan:
					// Graceful shutdown
					// nothing to close due to the socket closing after each loop
					log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
					return
				case <- loopPeriod:
					// Create the connection the server in every loop iteration
					c.createClientSocket()

					fmt.Fprintf(
						c.conn,
						"[CLIENT %v] Message N°%v\n",
						c.config.ID,
						msgID,
					)
					msg, err := bufio.NewReader(c.conn).ReadString('\n')
					c.conn.Close()
			
					if err != nil {
						log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
							c.config.ID,
							err,
						)
						return
					}
			
					log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
						c.config.ID,
						msg,
					)

					msgID ++
			}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}*/