package common

import (
	"net"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// BET_SIZE_SIZE Represents the size in bytes of the bet size field in the packet
const BET_SIZE_SIZE = 2
// AMOUNT_BETS_SIZE Represents the size in bytes of the amount of bets field in the packet
const AMOUNT_BETS_SIZE = 2
const CODE_SIZE = 1
const BET_BATCH_CODE = 0
const BET_RESPONSE_CODE = 1
const WINNER_REQUEST_CODE = 2
const WINNER_RESPONSE_CODE = 3

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
	err := binary.Write(&buf, binary.BigEndian, int8(BET_BATCH_CODE))
	if err != nil {
		return err
	}
	packet_size += 1

	// The amount of bets to be send is codified and added to buf
	err = binary.Write(&buf, binary.BigEndian, int16(len(bets)))
	if err != nil {
		return err
	}
	packet_size += AMOUNT_BETS_SIZE

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
		packet_size += BET_SIZE_SIZE

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

func (c *CentralLoteriaNacional) requestWinner(id int) error {
	var buf bytes.Buffer
	packet_size := 0

	err := binary.Write(&buf, binary.BigEndian, int8(WINNER_REQUEST_CODE))
	if err != nil {
		return err
	}
	packet_size += 1

	// The amount of bets to be send is codified and added to buf
	err = binary.Write(&buf, binary.BigEndian, int32(id))
	if err != nil {
		return err
	}
	packet_size += 4

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
	codeData := make([]byte, CODE_SIZE)
	for readBytes := 0; readBytes < CODE_SIZE; {
		n, err := c.conn.Read(codeData)
		if err != nil {
			return err
		}
		readBytes += n
	}

	var code int8
	err := binary.Read(bytes.NewReader(codeData), binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != BET_RESPONSE_CODE {
		message := fmt.Sprintf("Incorrect message received: Expecting message type %v, received message type %v", BET_RESPONSE_CODE, code)
		return errors.New(message)
	}

	data := make([]byte, AMOUNT_BETS_SIZE)
	for readBytes := 0; readBytes < AMOUNT_BETS_SIZE; {
		n, err := c.conn.Read(data)
		if err != nil {
			return err
		}
		readBytes += n
	}

	var bets_read int16
	err = binary.Read(bytes.NewReader(data), binary.BigEndian, &bets_read)
	if err != nil {
		return err
	}

	if bets_read != 0 {
		return nil
	}

	return errors.New("Server error while receiving the bet information")
}

func (c *CentralLoteriaNacional) readWinners() ([]int, error) {
	codeData := make([]byte, CODE_SIZE)
	for readBytes := 0; readBytes < CODE_SIZE; {
		n, err := c.conn.Read(codeData)
		if err != nil {
			return nil, err
		}
		readBytes += n
	}

	var code int8
	err := binary.Read(bytes.NewReader(codeData), binary.BigEndian, &code)
	if err != nil {
		return nil, err
	}

	if code != WINNER_RESPONSE_CODE {
		message := fmt.Sprintf("Incorrect message received: Expecting message type %v, received message type %v", WINNER_REQUEST_CODE, code)
		return nil, errors.New(message)
	}

	amountWinnersData := make([]byte, AMOUNT_BETS_SIZE)
	for readBytes := 0; readBytes < AMOUNT_BETS_SIZE; {
		n, err := c.conn.Read(amountWinnersData)
		if err != nil {
			return nil, err
		}
		readBytes += n
	}

	var amountWinners int16
	err = binary.Read(bytes.NewReader(amountWinnersData), binary.BigEndian, &amountWinners)
	if err != nil {
		return nil, err
	}

	winners := make([]int, amountWinners)
	for n := 0; n < int(amountWinners); n++ {
		documentData := make([]byte, 4)
		for readBytes := 0; readBytes < 4; {
			n, err := c.conn.Read(documentData)
			if err != nil {
				return nil, err
			}
			readBytes += n
		}

		var winner int32
		err := binary.Read(bytes.NewReader(documentData), binary.BigEndian, &winner)
		if err != nil {
			return nil, err
		}

		winners[n] = int(winner)
	}

	return winners, nil
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

func (c *CentralLoteriaNacional) GetWinners(id int) ([]int, error) {
	err := c.createSocket()
	if err != nil {
		return nil, err
	}

	err = c.requestWinner(id)
	if err != nil {
		return nil, err
	}

	winners, err := c.readWinners()
	if err != nil {
		return nil, err
	}

	c.conn.Close()

	return winners, nil
}