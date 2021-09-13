package main

import (
	"github.com/google/uuid"
)

// Email is an email entity
type Email struct {
	ID      uuid.UUID `json:"id"`
	To      []string  `json:"to"`
	CC      []string  `json:"cc"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	ReplyTo string    `json:"reply_to"`
	Body    string    `json:"body"`
}

// // EmailRepository handle the persistence of emails
// type EmailRepository interface {
// 	List(page, limit int) ([]Email, error)
// 	FindByID(id uuid.UUID) (Email, error)
// 	Save(email Email) error
// 	Delete(id uuid.UUID) error
// }

// NewEmailRepository instance
func NewEmailRepository(ds Datastore) *EmailRepository {
	return &EmailRepository{datastore: ds}
}

// EmailRepository stores and fetches items
type EmailRepository struct {
	datastore Datastore
}

// Get a single email
func (r *EmailRepository) Get(id uuid.UUID) (*Email, error) {
	var email *Email
	if err := r.datastore.Get(id, &email); err != nil {
		return nil, err
	}
	return email, nil
}

// Store a new email
func (r *EmailRepository) Store(email *Email) error {
	email.ID = uuid.New()
	return r.datastore.Store(email)
}

// List all emails
func (r *EmailRepository) List() ([]*Email, error) {
	var emails []*Email
	if err := r.datastore.List(&emails); err != nil {
		return nil, err
	}
	return emails, nil
}
