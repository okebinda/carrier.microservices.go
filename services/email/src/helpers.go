package main

import (
	"net/http"
	"strconv"
	"time"

	es "carrier.microservices.go/src/lib/email"
	"carrier.microservices.go/src/lib/store"
)

// GetQueryParamInt64 parses an int64 value from the URL query string, using a default if not present
func GetQueryParamInt64(r *http.Request, key string, def int64) (int64, error) {
	var value int64
	var err error

	if r.URL.Query().Has(key) {
		value, err = strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
	} else {
		value = def
	}
	return value, err
}

// SendEmail sends an email via the supplied service
func SendEmail(exchange es.EmailExchange, email *Email, emailRepository *EmailRepository) bool {

	sent := false

	// create change set for email
	changeSet := store.ChangeSet{
		"attempts": email.Attempts + 1,
	}

	// create email record to communicate with service
	exEmail := es.Email{
		Recipients:    email.Recipients,
		Template:      email.Template,
		Substitutions: email.Substitutions,
	}

	// send email and update record
	err := exchange.Send(&exEmail)
	if err != nil {
		logger.Errorf("Email exchange error: %s\n", err)
		changeSet["SendStatus"] = EmailStatusQueued
	} else {
		logger.Debugw("SparkPost transmission successful.")
		email.Queued = time.Time{}
		changeSet["send_status"] = EmailStatusComplete
		changeSet["service_id"] = exEmail.ID
		changeSet["last_attempt_at"] = exEmail.LastAttemptAt
		changeSet["accepted"] = exEmail.Accepted
		changeSet["rejected"] = exEmail.Rejected
		changeSet["queued"] = email.Queued
		sent = true
	}

	// save again with transmission data
	err = emailRepository.Update(email, changeSet)
	if err != nil {
		logger.Errorf("Unable to update email: %v", err)
	}

	return sent
}
