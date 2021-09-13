package validation

import (
	"reflect"
	"testing"
)

type TestPayload struct {
	Param1 string `json:"param1" validate:"required,min=2,max=10"`
	Param2 string `json:"param2" validate:"required,len=4"`
	Param3 int    `json:"param3" validate:"required,numeric"`
}

func TestCheckPass(t *testing.T) {
	payload := TestPayload{
		Param1: "testing",
		Param2: "ABCD",
		Param3: 25,
	}

	ok, _ := Check(payload)

	// test ok is true
	if !ok {
		t.Errorf("Check() returned false, expected true.")
	}
}

func TestCheckFail(t *testing.T) {
	expectedErrorMap := map[string]map[string]map[string]string{
		"errors": {
			"param1": {
				"min": "2",
			},
			"param2": {
				"len": "4",
			},
			"param3": {
				"required": "",
			},
		},
	}

	payload := TestPayload{
		Param1: "t",
		Param2: "ABCDE",
	}

	ok, errorMap := Check(payload)

	// test ok is false
	if ok {
		t.Errorf("Check() returned false, expected true.")
	}

	// test errorMap
	if !reflect.DeepEqual(errorMap, expectedErrorMap) {
		t.Errorf("errorMap was incorrect: got %v, expected %v.", errorMap, expectedErrorMap)
	}
}

type customInt struct {
	Value int
}

type TestPayload2 struct {
	Param1 string    `json:"param1" validate:"required,min=2,max=10"`
	Param2 string    `json:"param2" validate:"required,len=4"`
	Param3 int       `json:"param3" validate:"required,numeric"`
	Param4 customInt `json:"param4" validate:"required"`
}

// tests that AddCustomTypeFunc adds custom type functions to be regsitered
func TestAddCustomTypeFunc(t *testing.T) {
	testFunc := func(field reflect.Value) interface{} {
		return nil
	}

	// test length of customTypeFuncs slice before adding
	if len(customTypeFuncs) != 0 {
		t.Errorf("len(customTypeFuncs) was incorrect: got %d, expected %d.", len(customTypeFuncs), 0)
	}

	AddCustomTypeFunc(testFunc, customInt{})

	// test length of customTypeFuncs slice after adding
	if len(customTypeFuncs) != 1 {
		t.Errorf("len(customTypeFuncs) was incorrect: got %d, expected %d.", len(customTypeFuncs), 1)
	}

	// reset customTypeFuncs
	customTypeFuncs = customTypeFuncs[:0]
}

func TestCheckPassWithCustomTypeFunc(t *testing.T) {
	payload := TestPayload2{
		Param1: "testing",
		Param2: "ABCD",
		Param3: 25,
		Param4: customInt{5},
	}

	testFunc := func(field reflect.Value) interface{} {
		if ci, ok := field.Interface().(customInt); ok {
			emptyCI := customInt{}
			if ci == emptyCI {
				return false
			}
			return true
		}
		return nil
	}

	AddCustomTypeFunc(testFunc, customInt{})

	ok, _ := Check(payload)

	// test ok is true
	if !ok {
		t.Errorf("Check() returned false, expected true.")
	}

	// reset customTypeFuncs
	customTypeFuncs = customTypeFuncs[:0]
}

func TestCheckFailWithEmptyCustomTypeFunc(t *testing.T) {
	payload := TestPayload2{
		Param1: "testing",
		Param2: "ABCD",
		Param3: 25,
		Param4: customInt{},
	}

	testFunc := func(field reflect.Value) interface{} {
		if ci, ok := field.Interface().(customInt); ok {
			emptyCI := customInt{}
			if ci == emptyCI {
				return false
			}
			return true
		}
		return nil
	}

	AddCustomTypeFunc(testFunc, customInt{})

	ok, _ := Check(payload)

	// test ok is false
	if ok {
		t.Errorf("Check() returned true, expected false.")
	}

	// reset customTypeFuncs
	customTypeFuncs = customTypeFuncs[:0]
}
