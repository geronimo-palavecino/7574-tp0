package common

import (
	"net"
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

const MAX_AMOUNT_TRIES = 3
// BET_SIZE_SIZE Represents the size in bytes of the bet size field in the packet
const BET_SIZE_SIZE = 2
// RESPONSE_CODE_SIZE Represents the size in bytes of the response code from the server
const RESPONSE_CODE_SIZE = 2

// CentralLoteriaNacional Entity that encapsulates the communication with the server representing the Central de Loteria Nacional
type CentralLoteriaNacional struct {
	Address string
	conn 	net.Conn
}

// NewCentralLoteriaNacional Initializes a new CentralLoteriaNacional that
// will communicate with the given address
func NewCentralLoteriaNacional(address string) CentralLoteriaNacional {
	central := CentralLoteriaNacional{
		Address: address,
	}
	return central
}

// CreateSocket Initializes central socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *CentralLoteriaNacional) createSocket() error {
	var err error
	for tries := 1; tries <= MAX_AMOUNT_TRIES; tries ++ {
		conn, err := net.Dial("tcp", c.Address)
		if err == nil {
			c.conn = conn
			return err
		}
		
		time.Sleep(time.Duration(tries * 500) * time.Millisecond)
	}
	
	return err
}

// writeBet Sends the given Bet, in its byte representation
// though the underlying connection
func (c *CentralLoteriaNacional) writeBet(bet Bet) error {
	var buf bytes.Buffer
	packetSize := 0
	
	betBytes, err := bet.Bytes()
	if err != nil {
		return err
	}

	lenBetBytes := len(betBytes)
	err = binary.Write(&buf, binary.BigEndian, int16(lenBetBytes))
	if err != nil {
		return err
	}
	packetSize += BET_SIZE_SIZE

	buf.Write(betBytes)
	packetSize += lenBetBytes

	// Mechanism to avoid short write
	packet := buf.Bytes()
	for bytesSent := 0;  bytesSent < packetSize; {
		n, err := c.conn.Write(packet[bytesSent:])
		if err != nil {
			return err
		}
		bytesSent += n
	}
	
	return nil
}

// readConfirmation Reads the amount of bets the server read from a packet
func (c *CentralLoteriaNacional) readConfirmation() error {

	data := make([]byte, RESPONSE_CODE_SIZE)
	for readBytes := 0; readBytes < RESPONSE_CODE_SIZE; {
		n, err := c.conn.Read(data)
		if err != nil {
			return err
		}
		readBytes += n
	}
	var confirmation int16
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &confirmation)
	if err != nil {
		return err
	}

	if confirmation != 0 {
		return nil
	}

	return errors.New("Server error while receiving the bet information")
}

// sendBet Sends the given Bet through the underlying connection, and waits for
// confirmation of reception
func (c *CentralLoteriaNacional) SendBet(bet Bet) error {
	err := c.createSocket()
	if err != nil {
		return err
	}

	err = c.writeBet(bet)
	if err != nil {
		return err
	}

	err = c.readConfirmation()
	if err != nil {
		return err
	}

	c.conn.Close()

	return nil
}	