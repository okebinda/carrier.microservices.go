package main

import (
	"time"

	// "carrier.microservices.go/src/lib/datetime"
	"carrier.microservices.go/src/lib/store"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

// Update an existing email
func (r *EmailRepository) Update(email *Email) error {
	return r.datastore.Update(
		email.ID,
		map[string]*dynamodb.AttributeValue{
			":to": {
				SS: aws.StringSlice(email.To),
			},
			":cc": {
				SS: aws.StringSlice(email.CC),
			},
			":subject": {
				S: aws.String(email.Subject),
			},
			":from": {
				S: aws.String(email.From),
			},
			":reply_to": {
				S: aws.String(email.ReplyTo),
			},
			":body": {
				S: aws.String(email.Body),
			},
			// ":updated_at": {
			// 	// S: aws.String(datetime.JSONTime(time.Now())),
			// 	S: aws.String(time.Now().Format(datetime.ISO8601Datetime)),
			// },
		},
		// "set to_=:to, cc=:cc, subject=:subject, from_=:from, reply_to=:reply_to, body=:body, updated_at=:updated_at",
		"set to_=:to, cc=:cc, subject=:subject, from_=:from, reply_to=:reply_to, body=:body",
	)
}
