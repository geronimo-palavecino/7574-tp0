package common

import (
	"net"
	"bytes"
	"encoding/binary"
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

func (c *CentralLoteriaNacional) SendBet(bet Bet) error {
	err := c.createSocket()
	if err != nil {
		return err
	}

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

	// TODO: 
	// Enviar la longitud de la informacion (for ?)
	// hacer un for que no para hasta enviar todos los bytes de bet
	// esperar la respuesta del servidor
	// hacer un for que no para de leer hasta que se tomaron todos los bytes que indica el servidor
	
	return nil
}	