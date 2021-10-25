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
	SendStatus    int               `json:"send_status" validate:"numeric,gte=1,lte=4"`
	Queued        datetime.JSONTime `json:"queued"`
	// SendNow       *bool             `json:"send_now" validate:"required"`
	Priority  int    `json:"priority" validate:"required,numeric,gte=0,lte=3"`
	ServiceID string `json:"service_id"`
}

//BatchEmailRequestSchema defines the input shape and validation schema for
type BatchEmailRequestSchema struct {
	Emails []EmailRequestSchema `json:"emails" validate:"required,min=1"`
}

// EmailSchema defines the JSON schema for the Email model.
type EmailSchema struct {
	ID            uuid.UUID         `json:"id"`
	ServiceID     string            `json:"service_id"`
	Recipients    []string          `json:"recipients"`
	Template      string            `json:"template"`
	Substitutions map[string]string `json:"substitutions"`
	SendStatus    int               `json:"send_status"`
	Queued        datetime.JSONTime `json:"queued"`
	Priority      int               `json:"priority"`
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
	s.ServiceID = m.ServiceID
	s.Recipients = m.Recipients
	s.Template = m.Template
	s.Substitutions = m.Substitutions
	s.SendStatus = m.SendStatus
	s.Queued = datetime.JSONTime(m.Queued)
	s.Priority = m.Priority
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

// BatchEmailResponseSchema defines the response schema for a batch of Email records.
type BatchEmailResponseSchema struct {
	Emails []EmailSchema `json:"emails"`
	Sent   int64         `json:"sent"`
	Queued int64         `json:"queued"`
}
