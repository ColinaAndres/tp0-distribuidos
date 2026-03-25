package common

// BetConfig represents the configuration of a bet placed by a client
type BetConfig struct {
	Agency   string
	Name     string
	LastName string
	Document string
	Birth    string
	Number   string
}

// Bet represents a bet placed by a client
type Bet struct {
	agency   string
	name     string
	lastName string
	document string
	birth    string
	number   string
}

// NewBet Initializes a new bet receiving the configuration as a parameter
func NewBet(config BetConfig) *Bet {
	return &Bet{
		agency:   config.Agency,
		name:     config.Name,
		lastName: config.LastName,
		document: config.Document,
		birth:    config.Birth,
		number:   config.Number,
	}
}
