package fakes

import "net/http"

type ErrorWriter struct {
	Error error
}

func NewErrorWriter() *ErrorWriter {
	return &ErrorWriter{}
}

func (writer *ErrorWriter) Write(w http.ResponseWriter, err error) {
	writer.Error = err
}
