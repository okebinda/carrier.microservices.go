package datetime

import (
	"testing"
)

func TestISO8601Datetime(t *testing.T) {
	expectedFormat := "2006-01-02T15:04:05-0700"

	if ISO8601Datetime != expectedFormat {
		t.Errorf("ISO8601Datetime was incorrect: got %v, expected %v.", ISO8601Datetime, expectedFormat)
	}
}

func TestISO8601Date(t *testing.T) {
	expectedFormat := "2006-01-02"

	if ISO8601Date != expectedFormat {
		t.Errorf("ISO8601Date was incorrect: got %v, expected %v.", ISO8601Date, expectedFormat)
	}
}
