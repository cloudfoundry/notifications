package params

import "strings"

type ParseError struct{}

func (err ParseError) Error() string {
	return "Request body could not be parsed"
}

type SchemaError string

func NewSchemaError(msg string) SchemaError {
	return SchemaError(msg)
}

func (err SchemaError) Error() string {
	return string(err)
}

type ValidationError []string

func (err ValidationError) Error() string {
	return strings.Join(err, ", ")
}

func (err ValidationError) Errors() []string {
	return []string(err)
}
