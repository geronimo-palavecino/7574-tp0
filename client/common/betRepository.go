package common

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"time"
)

//FIRST_NAME_INDEX Index where the first name field of a bet is found in the csv
const FIRST_NAME_INDEX = 0
//LAST_NAME_INDEX Index where the last name field of a bet is found in the csv
const LAST_NAME_INDEX = 1
//DOCUMENT_INDEX Index where the document field of a bet is found in the csv
const DOCUMENT_INDEX = 2
//BIRTHDATE_INDEX Index where the birthdate field of a bet is found in the csv
const BIRTHDATE_INDEX = 3
//NUMBER_INDEX Index where the number field of a bet is found in the csv
const NUMBER_INDEX = 4

// BetRepository Entity that manages the bets
type BetRepository struct {
	filePath string
	file 	 *os.File
	scanner  *bufio.Scanner
}

// NewClient Initializes a new BetRepository receiving the filePath
// of the file to open as a parameter
func NewBetRepository(filePath string) (*BetRepository, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	betRepo := &BetRepository{
		filePath: filePath,
		file: file,
		scanner: scanner,
	}
	return betRepo, nil
}

func (r *BetRepository) GetBets(amount int, agency int) ([]Bet, error) {
	betsRead := 0 
	bets := []Bet{}
	
	for r.scanner.Scan() && betsRead < amount {
		betData := strings.Split(r.scanner.Text(), ",")

		birthdate, err := time.Parse("2006-01-02", betData[BIRTHDATE_INDEX])
		if err != nil {
			return nil, err
		}

		document, err := strconv.Atoi(betData[DOCUMENT_INDEX])
		if err != nil {
			return nil, err
		}

		number, err := strconv.Atoi(betData[NUMBER_INDEX])
		if err != nil {
			return nil, err
		}

		bet := Bet{
			Agency: agency,
			FirstName: betData[FIRST_NAME_INDEX],
			LastName: betData[LAST_NAME_INDEX],
			Document: document,
			Birthdate: birthdate,
			Number: number,
		}

		bets = append(bets, bet)
		betsRead += 1
	}

	err := r.scanner.Err()
	if err != nil {
		return nil, err
	}

	return bets, nil
}

func (r *BetRepository) Close() {
	r.file.Close()
}