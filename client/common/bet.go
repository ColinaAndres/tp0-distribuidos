package common

import (
	"os"
)

type Bet struct {
	agency   string
	name     string
	lastName string
	document string
	birth    string
	number   string
}

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
