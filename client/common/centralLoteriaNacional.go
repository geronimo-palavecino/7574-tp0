package common

import (
	"net"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const MAX_AMOUNT_TRIES = 3
// BET_SIZE_SIZE Represents the size in bytes of the bet size field in the packet
const BET_SIZE_SIZE = 2
// AMOUNT_BETS_SIZE Represents the size in bytes of the amount of bets field in the packet
const AMOUNT_BETS_SIZE = 2
// CODE_SIZE Represents the size in bytes of the code field in the packet
const CODE_SIZE = 1
// BET_BATCH_CODE Represents the code for a Bet Batch packet
const BET_BATCH_CODE = 0
// BET_RESPONSE_CODE Represents the code for a Bet Response packet
const BET_RESPONSE_CODE = 1
// WINNER_REQUEST_CODE Represents the code for a Winner Request packet
const WINNER_REQUEST_CODE = 2
// WINNER_RESPONSE_CODE Represents the code for a Winner Response packet
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
	var err error
	for tries := 1; tries <= MAX_AMOUNT_TRIES; tries++ {
		conn, err := net.Dial("tcp", c.Address)
		if err == nil {
			c.conn = conn
			return err
		}

		time.Sleep(time.Duration(tries * 500) * time.Millisecond)
	}
	
	return err
}

// writeBet Sends a Bets Batch packet with the given bets through the
// underlying socket
func (c *CentralLoteriaNacional) writeBet(bets []Bet) error {
	var buf bytes.Buffer

	// The packet code is codified and added to buf
	err := binary.Write(&buf, binary.BigEndian, int8(BET_BATCH_CODE))
	if err != nil {
		return err
	}

	// The amount of bets to be send is codified and added to buf
	err = binary.Write(&buf, binary.BigEndian, int16(len(bets)))
	if err != nil {
		return err
	}

	// All the bets are codified and added to the buf
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

		buf.Write(betBytes)
	}

	// The packet in the buf is sent
	err = writePacket(c.conn, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// requestWinner Sends a Winner Request packet indicating the agency whos
// asking through the underlying socket
func (c *CentralLoteriaNacional) requestWinner(id int) error {
	var buf bytes.Buffer

	// The packet code is codified and added to buf
	err := binary.Write(&buf, binary.BigEndian, int8(WINNER_REQUEST_CODE))
	if err != nil {
		return err
	}

	// The agency id is codified and added to buf
	err = binary.Write(&buf, binary.BigEndian, int32(id))
	if err != nil {
		return err
	}

	// The packet in the buf is sent
	err = writePacket(c.conn, buf.Bytes())
	if err != nil {
		return err
	}
	
	return nil
}

// readConfirmation Reads a Bet Response packet from the underlying socket
func (c *CentralLoteriaNacional) readConfirmation() error {
	// The packet code is read and checked to be correct
	code, err := readInt(c.conn, CODE_SIZE)
	if err != nil {
		return err
	}

	if code != BET_RESPONSE_CODE {
		message := fmt.Sprintf("Incorrect message received: Expecting message type %v, received message type %v", BET_RESPONSE_CODE, code)
		return errors.New(message)
	}

	// The number of bets is read and checked to be non zero
	betsRead, err := readInt(c.conn, AMOUNT_BETS_SIZE)
	if err != nil {
		return err
	}

	if betsRead != 0 {
		return nil
	}

	return errors.New("Server error while receiving the bet information")
}

// readWinners Reads a Winners Response packet from the underlying socket
func (c *CentralLoteriaNacional) readWinners() ([]int, error) {
	// The packet code is read and checked to be correct
	code, err := readInt(c.conn, CODE_SIZE)
	if err != nil {
		return nil, err
	}

	if code != WINNER_RESPONSE_CODE {
		message := fmt.Sprintf("Incorrect message received: Expecting message type %v, received message type %v", WINNER_REQUEST_CODE, code)
		return nil, errors.New(message)
	}

	// The number of winners is read
	amountWinners, err := readInt(c.conn, AMOUNT_BETS_SIZE)
	if err != nil {
		return nil, err
	}

	// The winners documents are read and added to the winners list
	winners := make([]int, amountWinners)
	for n := 0; n < int(amountWinners); n++ {
		winner, err := readInt(c.conn, 4)
		if err != nil {
			return nil, err
		}

		winners[n] = int(winner)
	}

	return winners, nil
}

// SendBet Sends the given Bets through the underlying connection, and waits for
// confirmation of reception
func (c *CentralLoteriaNacional) SendBets(bets []Bet) error {
	// The connection is opened
	err := c.createSocket()
	if err != nil {
		return err
	}

	// The Bet Batch packet is sent
	err = c.writeBet(bets)
	if err != nil {
		return err
	}

	// The Bet Confirmation packet is read
	err = c.readConfirmation()
	if err != nil {
		return err
	}

	// The connection is closed
	c.conn.Close()

	return nil
}	

// GetWinners Sends a request for the winers through the underlying connection,
// and waits for the response containing the documents of the winners from the agency
// with the given id
func (c *CentralLoteriaNacional) GetWinners(id int) ([]int, error) {
	// The connection is opened
	err := c.createSocket()
	if err != nil {
		return nil, err
	}

	// The Winners Request packet is sent
	err = c.requestWinner(id)
	if err != nil {
		return nil, err
	}

	// The Winners Response is read
	winners, err := c.readWinners()
	if err != nil {
		return nil, err
	}

	// The connection is closed
	c.conn.Close()

	return winners, nil
}

// writePacket Writes a packet in the given connection
func writePacket(conn net.Conn, packet []byte) error {
	// Mechanism to avoid short write
	for bytesSent := 0;  bytesSent < len(packet); {
		n, err := conn.Write(packet[bytesSent:])
		if err != nil {
			return err
		}
		bytesSent += n
	}

	return nil	
}

// readInt Reads an int constituted from n bytes from the given connection
func readInt(conn net.Conn, n int) (int32, error) {
	// The values data is read
	data := make([]byte, 4)
	for readBytes := 0; readBytes < n; {
		nRead, err := conn.Read(data[4-n-+readBytes:]) // This slicing method is done to add padding to the "leftside" of the data
		if err != nil {
			return 0, err
		}
		readBytes += nRead
	}

	// The value is converted to int32
	var value int32
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}

	return value, nil
}