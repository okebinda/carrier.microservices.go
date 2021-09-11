package main

import (
	"net/http"
)

// GetEmail retrieves a single emails
func GetEmail(w http.ResponseWriter, r *http.Request) {

	logger.Infow("GetEmail called")

	ctx := r.Context()
	email, ok := ctx.Value(keyEmail).(*Email)
	if !ok {
		userErrorResponse(w, 404, "Not found")
		return
	}

	// response
	successResponse(w, 200, map[string]interface{}{
		"email": email,
	})
}
