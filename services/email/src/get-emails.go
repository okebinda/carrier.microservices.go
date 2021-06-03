package main

import (
	"net/http"
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

	// response
	successResponse(w, 200, map[string]interface{}{
		"success": true,
	})
}
