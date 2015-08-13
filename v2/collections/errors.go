package collections

type PersistenceError struct {
	Err error
}

func (e PersistenceError) Error() string {
	return e.Err.Error()
}

type NotFoundError struct {
	Err error
}

func (e NotFoundError) Error() string {
	return e.Err.Error()
}

type DuplicateRecordError struct {
	Err error
}

func (e DuplicateRecordError) Error() string {
	return e.Err.Error()
}
