package valiant

import (
	"reflect"
	"strconv"
)

type RequiredFieldError struct {
	ErrorMessage string
}

func (err RequiredFieldError) Error() string {
	return err.ErrorMessage
}

func validateRequired(typed interface{}, untyped interface{}) error {
	parentValue, parentType := reflectedParent(typed)

	// TODO: Add JSON array type support here
	parentUntyped := untyped.(map[string]interface{})

	for i := 0; i < parentType.NumField(); i++ {
		childField := parentType.Field(i)
		childValue := parentValue.FieldByName(childField.Name)
		tag := childField.Tag.Get("json")

		err := validateRequiredField(childField, parentUntyped)
		if err != nil {
			return err
		}

		if childValue.Kind() == reflect.Struct {
			err = validateRequired(childValue.Interface(), parentUntyped[tag])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateRequiredField(field reflect.StructField, parentUntyped map[string]interface{}) error {
	if required(field) {
		if _, ok := parentUntyped[field.Tag.Get("json")]; !ok {
			return RequiredFieldError{ErrorMessage: "Missing required field '" + jsonName(field) + "'"}
		}
	}

	return nil
}

func required(field reflect.StructField) bool {
	if tag := field.Tag.Get("validate-required"); tag != "" {
		required, err := strconv.ParseBool(tag)
		if err != nil {
			panic("we got an error on parsing the validate-required tag")
		}

		return required
	}

	return false
}
