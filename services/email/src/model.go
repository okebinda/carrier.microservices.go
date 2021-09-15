package main

import (
	"time"

	"carrier.microservices.go/src/lib/store"
	"github.com/google/uuid"
)

// Email is an email entity
type Email struct {
	ID        uuid.UUID `json:"id"`
	To        []string  `json:"to_"`
	CC        []string  `json:"cc"`
	Subject   string    `json:"subject"`
	From      string    `json:"from_"`
	ReplyTo   string    `json:"reply_to"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EmailRepository stores and fetches items
type EmailRepository struct {
	datastore store.Datastore
}

// NewEmailRepository instance
func NewEmailRepository(ds store.Datastore) *EmailRepository {
	return &EmailRepository{datastore: ds}
}

// List all emails
func (r *EmailRepository) List() ([]*Email, error) {
	var emails []*Email
	if err := r.datastore.List(&emails); err != nil {
		return nil, err
	}
	return emails, nil
}

// Store a new email
func (r *EmailRepository) Store(email *Email) error {
	email.ID = uuid.New()
	email.CreatedAt = time.Now()
	email.UpdatedAt = time.Now()
	return r.datastore.Store(email)
}

// Get a single email
func (r *EmailRepository) Get(id uuid.UUID) (*Email, error) {
	var email *Email
	if err := r.datastore.Get(id, &email); err != nil {
		return nil, err
	}
	return email, nil
}

// Update an existing email
func (r *EmailRepository) Update(email *Email, changeSet store.ChangeSet) error {
	changeSet["updated_at"] = time.Now()
	return r.datastore.Update(email.ID, email, changeSet)
}
