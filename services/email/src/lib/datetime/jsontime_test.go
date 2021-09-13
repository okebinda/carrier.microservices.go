package datetime

import (
	"reflect"
	"testing"
	"time"
)

// tests the MarshalJSON() function of JSONTime
func TestMarshalJSON(t *testing.T) {
	expectedTime := []byte(`"2020-01-17T11:30:05+0000"`)

	jTime := JSONTime(time.Date(2020, time.Month(1), 17, 11, 30, 05, 0, time.UTC))
	b, err := jTime.MarshalJSON()

	// test no error
	if err != nil {
		t.Errorf("JSONTime.MarshalJSON() returned an error: got %v", err)
	}

	// test JSONTime.MarshalJSON() result
	if string(b) != string(expectedTime) {
		t.Errorf("JSONTime.MarshalJSON() was incorrect: got %v, expected %v", b, expectedTime)
	}
}

// tests the UnmarshalJSON() function of JSONTime
func TestUnmarshalJSON(t *testing.T) {
	expectedTime := time.Date(2020, time.Month(1), 17, 11, 30, 05, 0, time.UTC)

	jTime := JSONTime{}
	err := jTime.UnmarshalJSON([]byte(`"2020-01-17T11:30:05+0000"`))

	// test no error
	if err != nil {
		t.Errorf("JSONTime.UnmarshalJSON() returned an error: got %v", err)
	}

	// test JSONTime.UnmarshalJSON() result
	if !expectedTime.Equal(time.Time(jTime)) {
		t.Errorf("JSONTime.UnmarshalJSON() was incorrect: got %v, expected %v", time.Time(jTime), expectedTime)
	}
}

// tests the UnmarshalJSON() function of JSONTime with "null" value
func TestUnmarshalJSONNullValue(t *testing.T) {
	expectedTime := time.Date(1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)

	jTime := JSONTime{}
	err := jTime.UnmarshalJSON([]byte(`null`))

	// test no error
	if err != nil {
		t.Errorf("JSONTime.UnmarshalJSON() returned an error: got %v", err)
	}

	// test JSONTime.UnmarshalJSON() result
	if !expectedTime.Equal(time.Time(jTime)) {
		t.Errorf("JSONTime.UnmarshalJSON() was incorrect: got %v, expected %v", time.Time(jTime), expectedTime)
	}
}

// tests the ValidateJSONTime function returns true when JSONTime is not empty
func TestValidateJSONTimePass(t *testing.T) {
	jTime := JSONTime{}
	err := jTime.UnmarshalJSON([]byte(`"2020-01-17T11:30:05+0000"`))

	// test no error
	if err != nil {
		t.Errorf("JSONTime.UnmarshalJSON() returned an error: got %v", err)
	}

	// test pass
	if ValidateJSONTime(reflect.ValueOf(jTime)) != true {
		t.Errorf("JSONTime.ValidateJSONTime() did not pass when it should have")
	}
}

// tests the ValidateJSONTime function returns false when JSONTime is empty
func TestValidateJSONTimeFailure(t *testing.T) {
	jTime := JSONTime{}

	// test failure
	if ValidateJSONTime(reflect.ValueOf(jTime)) != false {
		t.Errorf("JSONTime.ValidateJSONTime() passed when it shouldn't have")
	}
}
