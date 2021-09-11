package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type key int

const (
	keyEmail key = iota
)

// Authorize checks if the request contains the proper authentication token
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check API key
		ok := authentication(r)
		if !ok {
			userErrorResponse(w, 403, "Permission denied.")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// authentication checks the request headers for an X-API-KEY value and compares it to env parameter
func authentication(r *http.Request) bool {
	APIKey := os.Getenv("API_KEY")
	if APIKey != "" {
		headerAPIKey := r.Header.Get("X-API-KEY")
		if headerAPIKey != APIKey {
			return false
		}
	}
	return true
}

// EmailCtx adds an Email object to the context if requested
func EmailCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get instance of emails table
		emailsTable := NewDynamoDBTable(db, "aws-com-kchevalier-dev-emails-table")

		// Create an email repository
		emailRepository := NewEmailRepository(emailsTable)

		// retrieve a single email
		id, err := uuid.Parse(chi.URLParam(r, "emailID"))
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

		ctx := context.WithValue(r.Context(), keyEmail, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
