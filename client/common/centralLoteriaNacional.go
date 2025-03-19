package common

import (
	"net"
	"bytes"
	"encoding/binary"
	"errors"
)

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
	conn, err := net.Dial("tcp", c.Address)
	c.conn = conn
	return err
}

// writeBet Sends the given Bet, in its byte representation
// though the underlying connection
func (c *CentralLoteriaNacional) writeBet(bet Bet) error {
	betBytes, err := bet.Bytes()
	if err != nil {
		return err
	}

	// The two added extra bytes are used to indicate the server the total
	//length of the message sent
	lenBetBytes := len(betBytes)

	var buf bytes.Buffer

	err = binary.Write(&buf, binary.BigEndian, int16(lenBetBytes))
	if err != nil {
		return err
	}

	message := append(buf.Bytes(), betBytes...)

	// Mechanism to avoid short write
	for bytesSent := 0;  bytesSent < lenBetBytes + 2; {
		n, err := c.conn.Write(message[bytesSent:])
		if err != nil {
			return err
		}
		bytesSent += n
	}
	
	return nil
}

// readConfirmation Reads the amount of bets the server read from a packet
func (c *CentralLoteriaNacional) readConfirmation() error {

	confirmation := make([]byte, 1)
	for readBytes := 0; readBytes < 1; {
		n, err := c.conn.Read(confirmation)
		if err != nil {
			return err
		}
		readBytes += n
	}

	if confirmation[0] != 0x00 {
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