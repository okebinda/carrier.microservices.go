package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ctf struct {
	Function validator.CustomTypeFunc
	Type     interface{}
}

var customTypeFuncs []*ctf

func AddCustomTypeFunc(function validator.CustomTypeFunc, customType interface{}) {
	customTypeFuncs = append(customTypeFuncs, &ctf{function, customType})
}

// Check performs validation on a struct using github.com/go-playground/validator rules.
func Check(s interface{}) (bool, map[string]map[string]map[string]string) {

	// init validation
	validate := validator.New()

	// register custom typ functions
	for _, ctf := range customTypeFuncs {
		validate.RegisterCustomTypeFunc(ctf.Function, ctf.Type)
	}

	// perform validation
	err := validate.Struct(s)

	// validation did not pass
	if err != nil {
		errorMap := make(map[string]map[string]map[string]string)
		errorMap["errors"] = map[string]map[string]string{}

		val := reflect.ValueOf(s)

		// loop over validation errors, add to output map
		for _, e := range err.(validator.ValidationErrors) {
			field, _ := val.Type().FieldByName(e.Field())
			jsonName := field.Tag.Get("json")
			if errorMap["errors"][jsonName] == nil {
				errorMap["errors"][jsonName] = map[string]string{}
			}
			errorMap["errors"][jsonName][e.Tag()] = e.Param()
		}

		return false, errorMap
	}

	// validation passed
	return true, nil
}
