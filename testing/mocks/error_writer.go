package mocks

import "net/http"

type ErrorWriter struct {
	WriteCall struct {
		Receives struct {
			Writer http.ResponseWriter
			Error  error
		}
	}
}

func NewErrorWriter() *ErrorWriter {
	return &ErrorWriter{}
}

func (ew *ErrorWriter) Write(writer http.ResponseWriter, err error) {
	ew.WriteCall.Receives.Writer = writer
	ew.WriteCall.Receives.Error = err
}
