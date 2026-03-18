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
		agency:   os.Getenv("AGENCY"),
		name:     os.Getenv("NAME"),
		lastName: os.Getenv("LAST_NAME"),
		document: os.Getenv("DOCUMENT"),
		birth:    os.Getenv("BIRTH"),
		number:   os.Getenv("NUMBER"),
	}
}
