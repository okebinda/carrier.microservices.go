package main

import (
	"net/http"
	"strconv"
)

// GetQueryParamInt64 parses an int64 value from the URL query string, using a default if not present
func GetQueryParamInt64(r *http.Request, key string, def int64) (int64, error) {
	var value int64
	var err error

	if r.URL.Query().Has(key) {
		value, err = strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
	} else {
		value = def
	}
	return value, err
}
