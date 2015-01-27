package valiant

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
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

	err = validateRequired(typed, untyped)
	if err != nil {
		return err
	}

	err = reportExtraKeys(typed, untyped)
	if err != nil {
		return err
	}

	return nil
}

func reflectedParent(typed interface{}) (reflect.Value, reflect.Type) {
	if reflect.TypeOf(typed).Kind() == reflect.Ptr {
		return reflect.ValueOf(typed).Elem(), reflect.TypeOf(typed).Elem()
	}

	return reflect.ValueOf(typed), reflect.TypeOf(typed)
}

func jsonName(field reflect.StructField) string {
	name := field.Tag.Get("json")
	name = strings.TrimSuffix(name, ",omitempty")

	if name != "" {
		return name
	}

	return field.Name
}
