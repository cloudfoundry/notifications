package models

type ErrRecordNotFound struct{}

func (err ErrRecordNotFound) Error() string {
	return "Record Not Found"
}

type ErrDuplicateRecord struct{}

func (err ErrDuplicateRecord) Error() string {
	return "Duplicate Record"
}

type TransactionCommitError struct {
	message string
}

func NewTransactionCommitError(msg string) TransactionCommitError {
	return TransactionCommitError{
		message: msg,
	}
}

func (err TransactionCommitError) Error() string {
	return err.message
}

type TemplateUpdateError struct {
	Message string
}

func (err TemplateUpdateError) Error() string {
	return err.Message
}
