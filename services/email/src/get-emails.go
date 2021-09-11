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
	// err = emailRepository.Store(&email)
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

	// response
	successResponse(w, 200, map[string]interface{}{
		"emails": emails,
	})
}
