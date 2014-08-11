package models

type ErrRecordNotFound struct{}

func (err ErrRecordNotFound) Error() string {
    return "Record Not Found"
}

type ErrDuplicateRecord struct{}

func (err ErrDuplicateRecord) Error() string {
    return "Duplicate Record"
}
