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
func (c *CentralLoteriaNacional) writeBet(bets []Bet) error {
	var buf bytes.Buffer
	packet_size := 0

	// The amount of bets to be send is codified and added to buf
	err := binary.Write(&buf, binary.BigEndian, int16(len(bets)))
	if err != nil {
		return err
	}
	packet_size += 2

	for n := 0; n < len(bets); n++ {
		betBytes, err := bets[n].Bytes()
		if err != nil {
			return err
		}

		lenBetBytes := len(betBytes)
		err = binary.Write(&buf, binary.BigEndian, int16(lenBetBytes))
		if err != nil {
			return err
		}
		packet_size += 2

		buf.Write(betBytes)
		packet_size += lenBetBytes
	}

	// Mechanism to avoid short write
	packet := buf.Bytes()
	for bytesSent := 0;  bytesSent < packet_size; {
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

	data := make([]byte, 2)
	for readBytes := 0; readBytes < 2; {
		n, err := c.conn.Read(data)
		if err != nil {
			return err
		}
		readBytes += n
	}

	var bets_read int16
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &bets_read)
	if err != nil {
		return err
	}

	if bets_read != 0 {
		return nil
	}

	return errors.New("Server error while receiving the bet information")
}

// sendBet Sends the given Bet through the underlying connection, and waits for
// confirmation of reception
func (c *CentralLoteriaNacional) SendBets(bets []Bet) error {
	err := c.createSocket()
	if err != nil {
		return err
	}

	err = c.writeBet(bets)
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