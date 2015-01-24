package rainmaker

import "fmt"

type UnauthorizedError struct {
	message []byte
}

func NewUnauthorizedError(message []byte) UnauthorizedError {
	return UnauthorizedError{
		message: message,
	}
}

func (err UnauthorizedError) Error() string {
	return fmt.Sprintf("Rainmaker UnauthorizedError: %s", err.message)
}

type NotFoundError struct {
	message []byte
}

func NewNotFoundError(message []byte) NotFoundError {
	return NotFoundError{
		message: message,
	}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("Rainmaker NotFoundError: %s", err.message)
}

type UnexpectedStatusError struct {
	status  int
	message []byte
}

func NewUnexpectedStatusError(status int, message []byte) UnexpectedStatusError {
	return UnexpectedStatusError{
		status:  status,
		message: message,
	}
}

func (err UnexpectedStatusError) Error() string {
	return fmt.Sprintf("Rainmaker UnexpectedStatusError: %d %s", err.status, err.message)
}

type ResponseReadError struct {
	internalError error
}

func NewResponseReadError(err error) ResponseReadError {
	return ResponseReadError{
		internalError: err,
	}
}

func (err ResponseReadError) Error() string {
	return "Rainmaker ResponseReadError: " + err.internalError.Error()
}

type ResponseBodyUnmarshalError struct {
	internalError error
}

func NewResponseBodyUnmarshalError(err error) ResponseBodyUnmarshalError {
	return ResponseBodyUnmarshalError{
		internalError: err,
	}
}

func (err ResponseBodyUnmarshalError) Error() string {
	return "Rainmaker ResponseBodyUnmarshalError: " + err.internalError.Error()
}

type RequestBodyMarshalError struct {
	internalError error
}

func NewRequestBodyMarshalError(err error) RequestBodyMarshalError {
	return RequestBodyMarshalError{internalError: err}
}

func (err RequestBodyMarshalError) Error() string {
	return "Rainmaker RequestBodyMarshalError: " + err.internalError.Error()
}

type RequestConfigurationError struct {
	internalError error
}

func NewRequestConfigurationError(err error) RequestConfigurationError {
	return RequestConfigurationError{internalError: err}
}

func (err RequestConfigurationError) Error() string {
	return "Rainmaker RequestConfigurationError: " + err.internalError.Error()
}

type RequestHTTPError struct {
	internalError error
}

func NewRequestHTTPError(err error) RequestHTTPError {
	return RequestHTTPError{internalError: err}
}

func (err RequestHTTPError) Error() string {
	return "Rainmaker RequestHTTPError: " + err.internalError.Error()
}
