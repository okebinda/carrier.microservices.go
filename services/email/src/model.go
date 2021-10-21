package main

import (
	"time"

	"carrier.microservices.go/src/lib/store"
	"github.com/google/uuid"
)

const (

	// EmailStatusQueued is a status constant for new/queued emails
	EmailStatusQueued = 1

	// EmailStatusProcessing is a status constant for emails that are attempting to send
	EmailStatusProcessing = 2

	// EmailStatusComplete is a status constant for emails that have been successfully sent
	EmailStatusComplete = 3

	// EmailStatusFailed is a status constant for emails that have failed to be sent and should not be tried again
	EmailStatusFailed = 4
)

// Email is an email entity
type Email struct {
	ID            uuid.UUID         `json:"id"`
	ServiceID     string            `json:"service_id"`
	Recipients    []string          `json:"recipients"`
	Template      string            `json:"template"`
	Substitutions map[string]string `json:"substitutions"`
	SendStatus    int               `json:"send_status"`
	Queued        time.Time         `json:"queued"`
	Attempts      int               `json:"attempts"`
	Accepted      int               `json:"accepted"`
	Rejected      int               `json:"rejected"`
	LastAttemptAt time.Time         `json:"last_attempt_at"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
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
func (r *EmailRepository) List(page, limit int64) ([]*Email, error) {
	var emails []*Email
	if err := r.datastore.List(&emails, page, limit); err != nil {
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

// Delete an existing email
func (r *EmailRepository) Delete(id uuid.UUID) error {
	return r.datastore.Delete(id)
}
