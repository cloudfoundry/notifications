package webutil

type ParseError struct{}

func (err ParseError) Error() string {
	return "Request body could not be parsed"
}

type SchemaError struct {
	Err error
}

func (e SchemaError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type MissingUserTokenError struct {
	Err error
}

func (e MissingUserTokenError) Error() string {
	return e.Err.Error()
}

type TemplateCreateError struct{}

func (err TemplateCreateError) Error() string {
	return "Failed to create Template in the database"
}
