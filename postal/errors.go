package postal

type UAA string

type UAAScopesError string

func (err UAAScopesError) Error() string {
	return string(err)
}

type UAAUserNotFoundError string

func (err UAAUserNotFoundError) Error() string {
	return string(err)
}

type TemplateLoadError string

func (err TemplateLoadError) Error() string {
	return string(err)
}

type CriticalNotificationError struct {
	kindID string
}

func NewCriticalNotificationError(kindID string) CriticalNotificationError {
	return CriticalNotificationError{kindID: kindID}
}

func (err CriticalNotificationError) Error() string {
	return "Insufficient privileges to send notification " + err.kindID
}
