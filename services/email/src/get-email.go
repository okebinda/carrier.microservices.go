package main

import (
	"net/http"
)

// GetEmail retrieves a single emails
func GetEmail(w http.ResponseWriter, r *http.Request) {

	logger.Debugw("GetEmail called")

	ctx := r.Context()
	email, ok := ctx.Value(keyEmail).(*Email)
	if !ok {
		logger.Errorf("Error retrieving email from context")
		serverErrorResponse(w)
	}

	// response
	successResponse(w, 200, map[string]interface{}{
		"email": email,
	})
}
