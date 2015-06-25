package services

type MissingKindOrClientError string

func (err MissingKindOrClientError) Error() string {
	return string(err)
}

type CriticalKindError string

func (err CriticalKindError) Error() string {
	return string(err)
}

type ClientMissingError string

func (err ClientMissingError) Error() string {
	return string(err)
}

type KindMissingError string

func (err KindMissingError) Error() string {
	return string(err)
}

type TemplateAssignmentError string

func (err TemplateAssignmentError) Error() string {
	return string(err)
}
