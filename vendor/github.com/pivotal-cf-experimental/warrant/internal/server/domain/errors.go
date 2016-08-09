package domain

type validationError string

func (e validationError) Error() string {
	return string(e)
}
