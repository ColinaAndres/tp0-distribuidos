package common

import (
	"os"
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
