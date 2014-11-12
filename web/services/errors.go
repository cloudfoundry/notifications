package services

type MissingKindOrClientError string

func (err MissingKindOrClientError) Error() string {
	return string(err)
}

type CriticalKindError string

func (err CriticalKindError) Error() string {
	return string(err)
}
