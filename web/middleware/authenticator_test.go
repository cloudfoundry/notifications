package middleware_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/web/middleware"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Authenticator", func() {
    var ware middleware.Authenticator
    var request *http.Request
    var writer *httptest.ResponseRecorder
    var token string

    BeforeEach(func() {
        var err error

        ware = middleware.NewAuthenticator()
        writer = httptest.NewRecorder()
        request, err = http.NewRequest("GET", "/some/path", nil)
        if err != nil {
            panic(err)
        }
    })

    Context("when the request uses an expired auth token", func() {
        BeforeEach(func() {
            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "mister-client",
                "cid":       "mister-client",
                "exp":       1404281214,
            }
            token = BuildToken(tokenHeader, tokenClaims)

            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)
        })

        It("returns a 401 status code and error message", func() {
            returnValue := ware.ServeHTTP(writer, request)

            Expect(returnValue).To(BeFalse())
            Expect(writer.Code).To(Equal(http.StatusUnauthorized))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: expired"))
        })
    })

    Context("when the request does not contain an auth token", func() {
        BeforeEach(func() {
            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
        })

        It("returns a 401 status code and error message", func() {
            returnValue := ware.ServeHTTP(writer, request)

            Expect(returnValue).To(BeFalse())
            Expect(writer.Code).To(Equal(http.StatusUnauthorized))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: missing"))
        })
    })

    Context("when the auth token does not contain the correct scope", func() {
        BeforeEach(func() {
            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "mister-client",
                "cid":       "mister-client",
                "exp":       3404281214,
                "scope":     []string{"cloud_controller.admin"},
            }
            token = BuildToken(tokenHeader, tokenClaims)

            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)
        })

        It("returns a 403 status code and error message", func() {
            returnValue := ware.ServeHTTP(writer, request)

            Expect(returnValue).To(BeFalse())
            Expect(writer.Code).To(Equal(http.StatusForbidden))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
        })
    })

    Context("when the auth token does not contain any scopes", func() {
        BeforeEach(func() {
            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "mister-client",
                "cid":       "mister-client",
                "exp":       3404281214,
            }
            token = BuildToken(tokenHeader, tokenClaims)

            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)
        })

        It("returns a 403 status code and error message", func() {
            returnValue := ware.ServeHTTP(writer, request)

            Expect(returnValue).To(BeFalse())
            Expect(writer.Code).To(Equal(http.StatusForbidden))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
        })
    })

    Context("when the request does not contain a auth valid token", func() {
        BeforeEach(func() {
            requestBody, err := json.Marshal(map[string]string{
                "kind": "forgot_password",
                "text": "Try to remember your password next time",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer something-invalid")
        })

        It("returns a 401 status code and error message", func() {
            returnValue := ware.ServeHTTP(writer, request)

            Expect(returnValue).To(BeFalse())
            Expect(writer.Code).To(Equal(http.StatusUnauthorized))

            parsed := map[string][]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: corrupt"))
        })
    })
})
