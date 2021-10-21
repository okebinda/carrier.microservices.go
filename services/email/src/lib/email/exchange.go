package mail

import (
	"time"
)

// Email represents and email to transmit
type Email struct {
	ID            string
	Recipients    []string
	Template      string
	Substitutions map[string]string
	Accepted      int
	Rejected      int
	LastAttemptAt time.Time
}

// Exchange is a generic interface for an email service
type Exchange interface {
	Init(ex *Exchange) error
	Send(email Email) error
}
