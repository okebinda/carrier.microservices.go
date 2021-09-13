package main

import (
	"net/http"
	// "github.com/google/uuid"
)

// GetEmails retrieves a list of emails
func GetEmails(w http.ResponseWriter, r *http.Request) {

	logger.Debugw("GetEmails called")

	// get instance of emails table
	emailsTable := NewDynamoDBTable(db, "aws-com-kchevalier-dev-emails-table")

	// Create an email repository
	emailRepository := NewEmailRepository(emailsTable)

	// // create a new test email
	// email := Email{
	// 	To:      []string{"test1@test.com", "test2@test.com"},
	// 	CC:      []string{"test3@test.com"},
	// 	Subject: "Hello",
	// 	From:    "test4@test.com",
	// 	ReplyTo: "test5@test.com",
	// 	Body:    "Lorem ipsum.",
	// }

	// // save email
	// err := emailRepository.Store(&email)
	// if err != nil {
	// 	logger.Errorf("Unable to save email: %v", err)
	// 	serverErrorResponse(w)
	// }

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
