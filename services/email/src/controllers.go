package main

import (
	"encoding/json"
	// "fmt"
	// "io"
	"net/http"
	// "github.com/google/uuid"
	"carrier.microservices.go/src/lib/validation"
)

// GetEmails retrieves a list of emails
func GetEmails(w http.ResponseWriter, r *http.Request) {

	logger.Debugw("GetEmails called")

	// get instance of emails table
	emailsTable := NewDynamoDBTable(db, "aws-com-kchevalier-dev-emails-table")

	// Create an email repository
	emailRepository := NewEmailRepository(emailsTable)

	// retrieve a list of emails
	emails, err := emailRepository.List()
	if err != nil {
		logger.Errorf("List emails error: %v", err)
		userErrorResponse(w, 404, "Not found")
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
		Page:   1,
		Limit:  10,
		Total:  1,
	})

	// // response
	// successResponse(w, 200, map[string]interface{}{
	// 	"emails": emails,
	// 	// "email": email,
	// })
}

// GetEmail retrieves a single emails
func GetEmail(w http.ResponseWriter, r *http.Request) {

	logger.Debugw("GetEmail called")

	// get email from context
	ctx := r.Context()
	email, ok := ctx.Value(keyEmail).(*Email)
	if !ok {
		logger.Errorf("Error retrieving email from context")
		serverErrorResponse(w)
	}

	// map result to response payload
	emailPayload := EmailSchema{}
	emailPayload.load(email)

	// response
	successResponse(w, 200, EmailResponseSchema{
		Email: emailPayload,
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

	// validate payload
	if ok, errorMap := validation.Check(payload); !ok {
		output, _ := json.Marshal(errorMap)
		generateResponse(w, http.StatusBadRequest, output)
		return
	}

	// get instance of email repository
	emailRepository := NewEmailRepository(NewDynamoDBTable(db, "aws-com-kchevalier-dev-emails-table"))

	// create a new email record
	email := Email{
		To:      payload.To,
		CC:      payload.CC,
		Subject: payload.Subject,
		From:    payload.From,
		ReplyTo: payload.ReplyTo,
		Body:    payload.Body,
	}

	// save email
	err = emailRepository.Store(&email)
	if err != nil {
		logger.Errorf("Unable to save email: %v", err)
		serverErrorResponse(w)
	}

	// map result to response payload
	emailPayload := EmailSchema{}
	emailPayload.load(&email)

	// response
	successResponse(w, 201, EmailResponseSchema{
		Email: emailPayload,
	})
}
