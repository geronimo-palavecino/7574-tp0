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

const MAX_AMOUNT_TRIES = 3

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
	config ClientConfig
	sigChan chan os.Signal
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	client := &Client{
		config: config,
		sigChan: sigChan,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	var err error
	for tries := 1; tries <= MAX_AMOUNT_TRIES; tries++ {
		conn, err := net.Dial("tcp", c.config.ServerAddress)
		if err == nil {
			c.conn = conn
			return nil
			
		}

		time.Sleep(time.Duration(tries * 500) * time.Millisecond)
	}
	
	return err
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.ClientLoop()
}

// ClientLoop Sen messages to the server until some time threshold is met or a SIGTERM is catched 
func (c *Client) ClientLoop() {
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
					err := c.createClientSocket()
					if err != nil {
						log.Criticalf(
							"action: connect | result: fail | client_id: %v | error: %v",
							c.config.ID,
							err,
						)
					}

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
}