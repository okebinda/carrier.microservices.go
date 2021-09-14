package main

import (
	"carrier.microservices.go/src/lib/datetime"
	"github.com/google/uuid"
)

// EmailRequestSchema defines the input validation schema for Email JSON requests.
type EmailRequestSchema struct {
	To      []string `json:"to" validate:"required,min=1,dive,required,email"`
	CC      []string `json:"cc" validate:"dive,email"`
	Subject string   `json:"subject" validate:"required,min=2,max=255"`
	From    string   `json:"from" validate:"required,email"`
	ReplyTo string   `json:"reply_to" validate:"email"`
	Body    string   `json:"body" validate:"required"`
}

// EmailSchema defines the JSON schema for the Email model.
type EmailSchema struct {
	ID        uuid.UUID         `json:"id"`
	To        []string          `json:"to"`
	CC        []string          `json:"cc"`
	Subject   string            `json:"subject"`
	From      string            `json:"from"`
	ReplyTo   string            `json:"reply_to"`
	Body      string            `json:"body"`
	CreatedAt datetime.JSONTime `json:"created_at"`
	UpdatedAt datetime.JSONTime `json:"updated_at"`
}

// Loads an Email record into EmailSchema.
func (s *EmailSchema) load(m *Email) {
	s.ID = m.ID
	s.To = m.To
	s.CC = m.CC
	s.Subject = m.Subject
	s.From = m.From
	s.ReplyTo = m.ReplyTo
	s.Body = m.Body
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
	Page   int           `json:"page"`
	Limit  int           `json:"limit"`
	Total  int           `json:"total"`
}
