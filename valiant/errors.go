package valiant

type ExtraFieldError struct {
	ErrorMessage string
}

func (err ExtraFieldError) Error() string {
	return err.ErrorMessage
}

type RequiredFieldError struct {
	ErrorMessage string
}

func (err RequiredFieldError) Error() string {
	return err.ErrorMessage
}
