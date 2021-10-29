package main

import (
	"context"
	"net/http"
	"os"

	emailService "carrier.microservices.go/src/lib/email"
	"carrier.microservices.go/src/lib/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type key int

const (
	keyEmail key = iota
	keyEmailRepository
	keyEmailExchange
)

// LogRequest logs the request
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check API key
		logger.Infow("Request",
			"Method", r.Method,
			"RequestURI", r.RequestURI,
			"RemoteAddr", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)
	})
}

// Authorize checks if the request contains the proper authentication token
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check API key
		ok := authentication(r)
		if !ok {
			userErrorResponse(w, 401, "Permission denied.")
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

// EmailRepositoryCtx adds a hepler function to the context to generate an instance of the EmailRepository
func EmailRepositoryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getEmailRepository := func() *EmailRepository {
			return NewEmailRepository(store.NewDynamoDBTable(db, os.Getenv("EMAILS_TABLE")))
		}
		ctx := context.WithValue(r.Context(), keyEmailRepository, getEmailRepository)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// EmailCtx adds an Email object to the context if requested
func EmailCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get email repository from context
		emailRepository := r.Context().Value(keyEmailRepository).(func() *EmailRepository)()

		// parse ID from URL into UUID
		id, err := uuid.Parse(chi.URLParam(r, "emailID"))
		if err != nil {
			logger.Errorf("Error creating UUID: %v", err)
			serverErrorResponse(w)
		}

		// retrieve a single email
		email, err := emailRepository.Get(id)
		if err != nil {
			switch err.(type) {
			case *store.NotFoundError:
				userErrorResponse(w, 404, "Not found")
			default:
				logger.Errorf("Unable to retrieve email from datastore: %v", err)
				serverErrorResponse(w)
			}
			return
		}

		ctx := context.WithValue(r.Context(), keyEmail, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// EmailExchangeCtx adds a hepler function to the context to generate an instance of the EmailRepository
func EmailExchangeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getEmailExchange := func() emailService.EmailExchange {
			return &emailService.SparkPostExchange{}
		}
		ctx := context.WithValue(r.Context(), keyEmailExchange, getEmailExchange)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
