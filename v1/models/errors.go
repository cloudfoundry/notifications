package models

type NotFoundError struct {
	Err error
}

func (e NotFoundError) Error() string {
	return e.Err.Error()
}

type DuplicateError struct {
	Err error
}

func (e DuplicateError) Error() string {
	return e.Err.Error()
}

type TransactionCommitError struct {
	Err error
}

func (e TransactionCommitError) Error() string {
	return e.Err.Error()
}

type TemplateUpdateError struct {
	Err error
}

func (e TemplateUpdateError) Error() string {
	return e.Err.Error()
}
