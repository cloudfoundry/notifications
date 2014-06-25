package stack

import (
    "bytes"
    "net/http"
    "net/http/httptest"
)

type Stack struct {
    Handler    http.Handler
    Middleware []Middleware
}

type Response struct {
    buffer *bytes.Buffer
    code   int
    header http.Header
}

func NewStack(handler http.Handler) Stack {
    return Stack{
        Handler: handler,
    }
}

func (s Stack) Use(wares ...Middleware) Stack {
    s.Middleware = append(s.Middleware, wares...)
    return s
}

func (s Stack) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    defer Recover(w, req)
    response := Response{
        buffer: bytes.NewBufferString(""),
        code:   http.StatusOK,
        header: make(http.Header),
    }

    var rec *httptest.ResponseRecorder

    for _, ware := range s.Middleware {
        rec = httptest.NewRecorder()
        halt := !ware.ServeHTTP(rec, req)
        s.Write(&response, rec)
        if halt {
            s.finalize(&response, w)
            return
        }
    }
    rec = httptest.NewRecorder()
    s.Handler.ServeHTTP(rec, req)
    s.Write(&response, rec)
    s.finalize(&response, w)
}

func (s *Stack) Write(response *Response, rec *httptest.ResponseRecorder) {
    for key, values := range rec.Header() {
        for _, value := range values {
            response.header.Add(key, value)
        }
    }
    response.buffer.Write(rec.Body.Bytes())
    response.code = rec.Code
}

func (s *Stack) finalize(response *Response, w http.ResponseWriter) {
    for key, values := range response.header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(response.code)
    w.Write(response.buffer.Bytes())
}
