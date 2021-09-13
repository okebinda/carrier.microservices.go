package main

import (
	"github.com/google/uuid"
	// "carrier.microservices.go/src/lib/datetime"
)

// // AppKeyAdminRequestSchema defines the input validation schema for admin AppKey JSON requests.
// type AppKeyAdminRequestSchema struct {
// 	Application string `json:"application" validate:"required,min=2,max=200"`
// 	Key         string `json:"key" validate:"required,len=32"`
// 	Status      int16  `json:"status" validate:"required,numeric,min=1,max=5"`
// }

// EmailSchema defines the JSON schema for the Email model.
type EmailSchema struct {
	ID      uuid.UUID `json:"id"`
	To      []string  `json:"to"`
	CC      []string  `json:"cc"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	ReplyTo string    `json:"reply_to"`
	Body    string    `json:"body"`
	// CreatedAt   datetime.JSONTime `json:"created_at"`
	// UpdatedAt   datetime.JSONTime `json:"updated_at"`
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
	// p.CreatedAt = datetime.JSONTime(m.CreatedAt)
	// p.UpdatedAt = datetime.JSONTime(m.UpdatedAt)
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
