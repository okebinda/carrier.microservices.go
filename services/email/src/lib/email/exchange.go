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

// EmailExchange is a generic interface for an email service
type EmailExchange interface {
	Init() error
	Send(email *Email) error
}
