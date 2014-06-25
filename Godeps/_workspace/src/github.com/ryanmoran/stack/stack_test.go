package stack_test

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "time"

    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("X-Handler", "my-handler")
    w.Write([]byte("Handler\n"))
}

type ChannelHandler struct {
    finish chan bool
    done   chan bool
}

func (h ChannelHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    <-h.finish
    h.done <- true
}

type Middleware struct {
    Name string
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) bool {
    str := fmt.Sprintf("Middleware: %s\n", m.Name)
    w.Header().Add("X-Middleware", m.Name)
    w.Write([]byte(str))
    return true
}

type HaltingMiddleware struct {
    Name string
}

func (m HaltingMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) bool {
    str := fmt.Sprintf("Middleware: %s\n", m.Name)
    w.WriteHeader(http.StatusUnauthorized)
    w.Write([]byte(str))
    return false
}

type PanickingMiddleware struct {
    err error
}

func (m PanickingMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) bool {
    if m.err != nil {
        panic(m.err)
    }
    return true
}

var _ = Describe("Stack", func() {
    It("assigns the Handler as the last object passed to the constructor", func() {
        handler := Handler{}
        s := stack.NewStack(handler).Use(Middleware{})
        Expect(s.Handler).To(Equal(handler))
    })

    It("assigns the remaining handlers as middleware", func() {
        middleware1 := Middleware{}
        middleware2 := Middleware{}
        s := stack.NewStack(Handler{}).Use(middleware1, middleware2)
        Expect(len(s.Middleware)).To(Equal(2))
        Expect(s.Middleware).To(ContainElement(middleware1))
        Expect(s.Middleware).To(ContainElement(middleware2))
    })

    It("executes the stack in order", func() {
        s := stack.NewStack(Handler{}).Use(Middleware{"first"}, Middleware{"second"})
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        s.ServeHTTP(writer, request)
        Expect(strings.Split(writer.Body.String(), "\n")).To(Equal([]string{
            "Middleware: first",
            "Middleware: second",
            "Handler",
            "",
        }))
    })

    It("stops cascading the handlers if a middleware handler returns false", func() {
        s := stack.NewStack(Handler{}).Use(Middleware{"first"}, HaltingMiddleware{"second"})
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        s.ServeHTTP(writer, request)
        Expect(strings.Split(writer.Body.String(), "\n")).To(Equal([]string{
            "Middleware: first",
            "Middleware: second",
            "",
        }))
        Expect(writer.Code).To(Equal(http.StatusUnauthorized))
    })

    It("clears the buffer after the stack completes", func() {
        s := stack.NewStack(Handler{}).Use(Middleware{"first"}, Middleware{"second"})
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        s.ServeHTTP(writer, request)
        Expect(strings.Split(writer.Body.String(), "\n")).To(Equal([]string{
            "Middleware: first",
            "Middleware: second",
            "Handler",
            "",
        }))

        writer = httptest.NewRecorder()
        s.ServeHTTP(writer, request)
        Expect(strings.Split(writer.Body.String(), "\n")).To(Equal([]string{
            "Middleware: first",
            "Middleware: second",
            "Handler",
            "",
        }))
    })

    It("records and returns headers", func() {
        s := stack.NewStack(Handler{}).Use(Middleware{"first"}, Middleware{"second"})
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        s.ServeHTTP(writer, request)
        Expect(writer.Header()["X-Middleware"]).To(Equal([]string{"first", "second"}))
        Expect(writer.Header()["X-Handler"]).To(Equal([]string{"my-handler"}))
    })

    Describe("handles concurrent requests", func() {
        It("ensures the headers, code, and body of a response are unique to a request", func() {
            finish := make(chan bool)
            done := make(chan bool)

            s := stack.NewStack(ChannelHandler{finish, done}).Use(Middleware{"concurrent"})
            writer1 := httptest.NewRecorder()
            request1, err := http.NewRequest("GET", "/my/path", nil)
            if err != nil {
                panic(err)
            }

            writer2 := httptest.NewRecorder()
            request2, err := http.NewRequest("GET", "/my/path", nil)
            if err != nil {
                panic(err)
            }

            go func() {
                s.ServeHTTP(writer1, request1)
            }()

            go func() {
                s.ServeHTTP(writer2, request2)
            }()

            finish <- true
            finish <- true
            <-done
            <-done

            <-time.After(10 * time.Millisecond)

            header1 := writer1.Header()["X-Middleware"]
            header2 := writer2.Header()["X-Middleware"]

            Expect(len(header1)).To(Equal(1))
            Expect(len(header2)).To(Equal(1))
        })
    })
})
