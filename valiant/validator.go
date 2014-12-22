package valiant

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
)

type Validator struct {
	input io.Reader
}

func NewValidator(input io.Reader) Validator {
	return Validator{
		input: input,
	}
}

func (validator Validator) Validate(typed interface{}) error {
	typedInput := bytes.NewBuffer([]byte{})
	untypedInput := bytes.NewBuffer([]byte{})
	_, err := io.Copy(io.MultiWriter(typedInput, untypedInput), validator.input)
	if err != nil {
		return err
	}

	err = json.NewDecoder(typedInput).Decode(typed)
	if err != nil {
		return err
	}

	var untyped interface{}
	err = json.NewDecoder(untypedInput).Decode(&untyped)
	if err != nil {
		return err
	}

	return validate(typed, untyped)
}

func reflectedParent(typed interface{}) (reflect.Value, reflect.Type) {
	if reflect.TypeOf(typed).Kind() == reflect.Ptr {
		return reflect.ValueOf(typed).Elem(), reflect.TypeOf(typed).Elem()
	}

	return reflect.ValueOf(typed), reflect.TypeOf(typed)
}

func validate(typed interface{}, untyped interface{}) error {
	parentValue, parentType := reflectedParent(typed)

	// TODO: Add JSON array type support here
	parentUntyped := untyped.(map[string]interface{})

	for i := 0; i < parentType.NumField(); i++ {
		childField := parentType.Field(i)
		childValue := parentValue.FieldByName(childField.Name)

		err := validateField(childField, parentUntyped)
		if err != nil {
			return err
		}

		if childValue.Kind() == reflect.Struct {
			err = validate(childValue.Interface(), parentUntyped[childField.Tag.Get("json")])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateField(field reflect.StructField, parentUntyped map[string]interface{}) error {
	if required(field) {
		if _, ok := parentUntyped[field.Tag.Get("json")]; !ok {
			return errors.New("Missing required field " + field.Name)
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
