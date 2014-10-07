package fakes

import "net/http"

type FakeErrorWriter struct {
    Error error
}

func NewFakeErrorWriter() *FakeErrorWriter {
    return &FakeErrorWriter{}
}

func (writer *FakeErrorWriter) Write(w http.ResponseWriter, err error) {
    writer.Error = err
}
