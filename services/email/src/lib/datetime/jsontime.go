package datetime

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"carrier.microservices.go/src/lib/validation"
)

// JSONTime is a formatted Time type for JSON
type JSONTime time.Time

func init() {

	// add custom type validation to validator
	validation.AddCustomTypeFunc(ValidateJSONTime, JSONTime{})
}

// MarshalJSON marshals a timestamp into formatted JSON
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(ISO8601Datetime))
	return []byte(stamp), nil
}

// UnmarshalJSON unmarshals a formatted JSON timestamp into a timestamp
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t = &JSONTime{}
		return nil
	}
	stamp, err := time.Parse(ISO8601Datetime, s)
	if err != nil {
		return err
	}
	*t = JSONTime(stamp)
	return nil
}

// ValidateJSONTime validates that JSONTime is not empty
func ValidateJSONTime(field reflect.Value) interface{} {
	if jsonTime, ok := field.Interface().(JSONTime); ok {
		emptyJSONTime := JSONTime{}
		if jsonTime == emptyJSONTime {
			return false
		}
		return true
	}
	return nil
}
