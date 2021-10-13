package main

import (
	"carrier.microservices.go/src/lib/datetime"
	"github.com/google/uuid"
)

// EmailRequestSchema defines the input validation schema for Email JSON requests.
type EmailRequestSchema struct {
	Recipients    []string          `json:"recipients" validate:"required,min=1,dive,required,email"`
	Template      string            `json:"template" validate:"required,min=2,max=255"`
	Substitutions map[string]string `json:"substitutions"`
	SendStatus    int               `json:"send_status" validate:"required,numeric,gte=1,lte=3"`
	Queued        datetime.JSONTime `json:"queued"`
}

// EmailSchema defines the JSON schema for the Email model.
type EmailSchema struct {
	ID            uuid.UUID         `json:"id"`
	Recipients    []string          `json:"recipients"`
	Template      string            `json:"template"`
	Substitutions map[string]string `json:"substitutions"`
	SendStatus    int               `json:"send_status"`
	Queued        datetime.JSONTime `json:"queued"`
	Attempts      int               `json:"attempts"`
	Accepted      int               `json:"accepted"`
	Rejected      int               `json:"rejected"`
	LastAttemptAt datetime.JSONTime `json:"last_attempt_at"`
	CreatedAt     datetime.JSONTime `json:"created_at"`
	UpdatedAt     datetime.JSONTime `json:"updated_at"`
}

// Loads an Email record into EmailSchema.
func (s *EmailSchema) load(m *Email) {
	s.ID = m.ID
	s.Recipients = m.Recipients
	s.Template = m.Template
	s.Substitutions = m.Substitutions
	s.SendStatus = m.SendStatus
	s.Queued = datetime.JSONTime(m.Queued)
	s.Attempts = m.Attempts
	s.Accepted = m.Accepted
	s.Rejected = m.Rejected
	s.LastAttemptAt = datetime.JSONTime(m.LastAttemptAt)
	s.CreatedAt = datetime.JSONTime(m.CreatedAt)
	s.UpdatedAt = datetime.JSONTime(m.UpdatedAt)
}

// EmailResponseSchema defines the response schema for a single Email record.
type EmailResponseSchema struct {
	Email EmailSchema `json:"email"`
}

// EmailListResponseSchema defines the response schema for a list of Email records.
type EmailListResponseSchema struct {
	Emails []EmailSchema `json:"emails"`
	Page   int64         `json:"page"`
	Limit  int64         `json:"limit"`
}
