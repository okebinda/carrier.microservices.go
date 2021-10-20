package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"carrier.microservices.go/src/lib/store"
	"carrier.microservices.go/src/lib/validation"
	sp "github.com/SparkPost/gosparkpost"
)

// GetEmails retrieves a list of emails
func GetEmails(w http.ResponseWriter, r *http.Request) {
	var page, limit int64
	var err error

	logger.Debugw("GetEmails called")

	// get page from query string
	page, err = GetQueryParamInt64(r, "page", 1)
	if err != nil || page < 1 {
		userErrorResponse(w, http.StatusBadRequest, "Invalid value for query parameter: page")
		return
	}

	// get limit from query string
	limit, err = GetQueryParamInt64(r, "limit", 25)
	if err != nil || limit < 1 || limit > 200 {
		userErrorResponse(w, http.StatusBadRequest, "Invalid value for query parameter: limit")
		return
	}

	// get email repository from context
	emailRepository := r.Context().Value(keyEmailRepository).(func() *EmailRepository)()

	// retrieve a list of emails
	emails, err := emailRepository.List(page, limit)
	if err != nil {
		logger.Errorf("List emails error: %v", err)
		userErrorResponse(w, 404, "Not found")
		return
	}

	// map results to response payload
	emailsPayload := []EmailSchema{}
	for _, email := range emails {
		emailPayload := EmailSchema{}
		emailPayload.load(email)
		emailsPayload = append(emailsPayload, emailPayload)
	}

	// response
	successResponse(w, 200, EmailListResponseSchema{
		Emails: emailsPayload,
		Page:   page,
		Limit:  limit,
	})
}

// PostEmails creates a new email record
func PostEmails(w http.ResponseWriter, r *http.Request) {
	var payload EmailRequestSchema
	var err error

	logger.Debugw("PostEmails called")

	// get payload from request body
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&payload); err != nil {
		userErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	logger.Debugf("Request payload: %+v", payload)

	// validate payload
	if ok, errorMap := validation.Check(payload); !ok {
		output, _ := json.Marshal(errorMap)
		generateResponse(w, http.StatusBadRequest, output)
		return
	}

	// get email repository from context
	emailRepository := r.Context().Value(keyEmailRepository).(func() *EmailRepository)()

	// create a new email record
	email := Email{
		Recipients:    payload.Recipients,
		Template:      payload.Template,
		Substitutions: payload.Substitutions,
		SendStatus:    payload.SendStatus,
		Queued:        time.Now(),
	}

	// set status differently if sending email now or later
	if *payload.SendNow == true {
		email.SendStatus = EmailStatusProcessing
	} else {
		email.SendStatus = EmailStatusQueued
	}

	// save email
	err = emailRepository.Store(&email)
	if err != nil {
		logger.Errorf("Unable to save email: %v", err)
		serverErrorResponse(w)
		return
	}

	// send email now
	if *payload.SendNow == true {

		logger.Debugw("Attempt to send email synchronously")

		// create change set for email
		changeSet := store.ChangeSet{
			"attempts":        email.Attempts + 1,
			"last_attempt_at": time.Now(),
		}

		// get SparkPost configuration from ENV
		apiKey := os.Getenv("SPARKPOST_API_KEY")
		apiBaseURL := os.Getenv("SPARKPOST_BASE_URL")
		apiVersion, _ := strconv.Atoi(os.Getenv("SPARKPOST_API_VERSION"))

		// create SparkPost client
		cfg := &sp.Config{
			BaseUrl:    apiBaseURL,
			ApiKey:     apiKey,
			ApiVersion: apiVersion,
		}
		var client sp.Client
		err := client.Init(cfg)
		if err != nil {
			logger.Errorf("SparkPost client init failed: %s\n", err)
		}

		// create recipient list
		recipients := []sp.Recipient{}
		for _, address := range email.Recipients {
			recipient := sp.Recipient{
				Address:          address,
				SubstitutionData: email.Substitutions,
			}
			recipients = append(recipients, recipient)
		}

		// send email
		tx := &sp.Transmission{
			Recipients: recipients,
			Content: map[string]interface{}{
				"template_id": email.Template,
			},
		}
		id, res, err := client.Send(tx)
		if err != nil {
			logger.Errorf("SparkPost email transmission failed: %s\n", err)
			changeSet["SendStatus"] = EmailStatusQueued
		} else {
			logger.Debugf("SparkPost transmission successful. ID: %v, Results: %v", id, res.Results)
			txResults := res.Results.(map[string]interface{})
			changeSet["send_status"] = EmailStatusComplete
			changeSet["accepted"] = int(txResults["total_accepted_recipients"].(float64))
			changeSet["rejected"] = int(txResults["total_rejected_recipients"].(float64))
			changeSet["queued"] = time.Time{}
			email.Queued = time.Time{}
		}

		// save again with transmission data
		err = emailRepository.Update(&email, changeSet)
		if err != nil {
			logger.Errorf("Unable to update email: %v", err)
			serverErrorResponse(w)
			return
		}
	}

	// map result to response payload
	emailPayload := EmailSchema{}
	emailPayload.load(&email)

	// response
	successResponse(w, 201, EmailResponseSchema{
		Email: emailPayload,
	})
}

// GetEmail retrieves a single email
func GetEmail(w http.ResponseWriter, r *http.Request) {

	logger.Debugw("GetEmail called")

	// get email from context
	email := r.Context().Value(keyEmail).(*Email)

	logger.Debugf("Email: %v", email)

	// map result to response payload
	emailPayload := EmailSchema{}
	emailPayload.load(email)

	// response
	successResponse(w, 200, EmailResponseSchema{
		Email: emailPayload,
	})
}

// UpdateEmail updates a single email
func UpdateEmail(w http.ResponseWriter, r *http.Request) {
	var payload EmailRequestSchema
	var err error

	logger.Debugw("UpdateEmail called")

	// get email from context
	ctx := r.Context()
	email := ctx.Value(keyEmail).(*Email)

	logger.Debugf("Email (before): %v", email)

	// get payload from request body
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&payload); err != nil {
		userErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// validate payload
	if ok, errorMap := validation.Check(payload); !ok {
		output, _ := json.Marshal(errorMap)
		generateResponse(w, http.StatusBadRequest, output)
		return
	}

	// get email repository from context
	emailRepository := ctx.Value(keyEmailRepository).(func() *EmailRepository)()

	// create change set for email
	changeSet := store.ChangeSet{
		"recipients":    payload.Recipients,
		"template":      payload.Template,
		"substitutions": payload.Substitutions,
		"send_status":   payload.SendStatus,
		"queued":        time.Time(payload.Queued),
	}

	// save email
	err = emailRepository.Update(email, changeSet)
	if err != nil {
		logger.Errorf("Unable to update email: %v", err)
		serverErrorResponse(w)
		return
	}

	// fix empty `queued` in result payload
	if time.Time(payload.Queued).IsZero() {
		email.Queued = time.Time{}
	}

	logger.Debugf("Email (after): %v", email)

	// map result to response payload
	emailPayload := EmailSchema{}
	emailPayload.load(email)

	// response
	successResponse(w, 200, EmailResponseSchema{
		Email: emailPayload,
	})
}

// DeleteEmail deletes a single email
func DeleteEmail(w http.ResponseWriter, r *http.Request) {
	var err error

	logger.Debugw("DeleteEmail called")

	// get email from context
	ctx := r.Context()
	email := ctx.Value(keyEmail).(*Email)

	logger.Debugf("Email: %v", email)

	// get email repository from context
	emailRepository := ctx.Value(keyEmailRepository).(func() *EmailRepository)()

	// delete email
	err = emailRepository.Delete(email.ID)
	if err != nil {
		logger.Errorf("Unable to delete email: %v", err)
		serverErrorResponse(w)
		return
	}

	// response
	successResponse(w, 204, nil)
}
