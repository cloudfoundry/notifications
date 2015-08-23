package mocks

import "net/http"

type ErrorWriter struct {
	WriteCall struct {
		Receives struct {
			Error error
		}
	}
}

func NewErrorWriter() *ErrorWriter {
	return &ErrorWriter{}
}

func (writer *ErrorWriter) Write(w http.ResponseWriter, err error) {
	writer.WriteCall.Receives.Error = err
	w.WriteHeader(http.StatusTeapot)
}
