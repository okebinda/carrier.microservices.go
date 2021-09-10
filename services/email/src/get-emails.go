package main

import (
	"net/http"
	"os"

	"github.com/google/uuid"
)

// GetEmails retrieves a lst of emails
func GetEmails(w http.ResponseWriter, r *http.Request) {

	logger.Infow("GetEmails called")

	// check API key
	ok := authentication(r)
	if !ok {
		userErrorResponse(w, 403, "Permission denied.")
		return
	}

	// connect to datastore
	conn, err := CreateConnection(os.Getenv("DYNAMODB_ENDPOINT"))
	if err != nil {
		logger.Errorf("Unable to connect to database: %v", err)
		serverErrorResponse(w)
	}

	// get instance of emails table
	emailsTable := NewDynamoDB(conn, "aws-com-kchevalier-dev-emails-table")

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

	// // retrieve a list of emails
	// emails, err := emailRepository.List()
	// if err != nil {
	// 	logger.Errorf("Unable to list emails: %v", err)
	// 	userErrorResponse(w, 404, "Not found")
	// }

	// retrieve a single email
	// id, err := uuid.FromBytes([]byte("ef7f232d-e100-4552-9b80-a6fd587ade36"))
	id, err := uuid.Parse("ef7f232d-e100-4552-9b80-a6fd587ade36")
	if err != nil {
		logger.Errorf("Error creating UUID: %v", err)
		serverErrorResponse(w)
	}

	email, err := emailRepository.Get(id)
	if err != nil {
		logger.Errorf("Unable to find email: %v", err)
		userErrorResponse(w, 404, "Not found")
		return
	}

	// response
	successResponse(w, 200, map[string]interface{}{
		"success": true,
		"email":   email,
	})
}
