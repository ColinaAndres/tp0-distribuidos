package common

import (
	"os"
	"strings"
)

const (
	name = iota
	lastName
	document
	birth
	number
)

// Bet represents a bet placed by a client
type Bet struct {
	agency   string
	name     string
	lastName string
	document string
	birth    string
	number   string
}

// NewBetFromLine creates a new Bet instance by parsing a line of text
// The line is expected to be in the format: "name,lastName,document,birth,number"
func NewBetFromLine(agency string, line string) *Bet {
	betFields := strings.Split(line, ",")
	return &Bet{
		agency:   agency,
		name:     betFields[name],
		lastName: betFields[lastName],
		document: betFields[document],
		birth:    betFields[birth],
		number:   betFields[number],
	}
}

// NewBetFromEnv creates a new Bet instance by reading
// the necessary fields from environment variables.
func NewBetFromEnv() *Bet {
	return &Bet{
		agency:   os.Getenv("CLI_ID"),
		name:     os.Getenv("NOMBRE"),
		lastName: os.Getenv("APELLIDO"),
		document: os.Getenv("DOCUMENTO"),
		birth:    os.Getenv("NACIMIENTO"),
		number:   os.Getenv("NUMERO"),
	}
}
