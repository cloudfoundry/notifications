package valiant

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

func ValidateJSON(typedData interface{}, rawJSON []byte) error {
	var untypedData interface{}
	err := json.Unmarshal(rawJSON, &untypedData)
	if err != nil {
		panic(err)
	}

	return recursivelyValidate(typedData, untypedData)
}

func recursivelyValidate(typedData interface{}, untypedData interface{}) error {
	parentType := reflect.TypeOf(typedData)
	parentValue := reflect.ValueOf(typedData)

	// Add JSON array type support here
	parentUntypedData := untypedData.(map[string]interface{})

	for i := 0; i < parentType.NumField(); i++ {
		childField := parentType.Field(i)

		err := validateSingleField(childField, parentUntypedData)
		if err != nil {
			return err
		}

		childValue := parentValue.FieldByName(childField.Name)
		jsonFieldName := childField.Tag.Get("json")

		if childValue.Kind() == reflect.Struct {
			err = recursivelyValidate(childValue.Interface(), parentUntypedData[jsonFieldName])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validateSingleField(field reflect.StructField, parentUntypedData map[string]interface{}) error {
	jsonFieldName := field.Tag.Get("json")

	var err error
	var required bool

	if tagText := field.Tag.Get("validate-required"); tagText != "" {
		required, err = strconv.ParseBool(tagText)
		if err != nil {
			panic("we got an error on parsing the validate-required tag")
		}
	}

	if required {
		if _, ok := parentUntypedData[jsonFieldName]; !ok {
			return errors.New("Missing required field " + field.Name)
		}
	}
	return nil
}
