package common

import (
	"time"
	"bytes"
	"encoding/binary"
)

// Bet Bet information
type Bet struct {
	Agency 		int
	FirstName 	string
	LastName 	string
	Document 	int
	Birthdate 	time.Time
	Number 		int
}

// Bytes Encodes a Bet into its bytes representation
func (b *Bet) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	// Encoding the Agency number into bytes
	err := binary.Write(&buf, binary.BigEndian, int32(b.Agency))
	if err != nil {
		return nil, err
	}

	// Encoding the First Name into bytes
	lenFirstName := uint8(len(b.FirstName))
	err = binary.Write(&buf, binary.BigEndian, lenFirstName)
	if err != nil {
		return nil, err
	}
	buf.WriteString(b.FirstName)

	// Encoding the Last Name into bytes
	lenLastName := uint8(len(b.LastName))
	err = binary.Write(&buf, binary.BigEndian, lenLastName)
	if err != nil {
		return nil, err
	}
	buf.WriteString(b.LastName)

	// Encoding the Document number into bytes
	err = binary.Write(&buf, binary.BigEndian, int32(b.Document))
	if err != nil {
		return nil, err
	}

	// Encoding the Birthdate into bytes
	buf.WriteString(b.Birthdate.Format("2006-01-02"))

	// Encoding the played Number into bytes
	err = binary.Write(&buf, binary.BigEndian, int32(b.Number))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}