package handlers_test

import (
    "net/http"
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUserParams", func() {
    Describe("ParseRequestBody", func() {
        It("parses the body of the given request", func() {
            request, err := http.NewRequest("GET", "/users", strings.NewReader(`{
                "kind": "test_email",
                "kind_description": "Descriptive Email Name",
                "source_description": "Descriptive Component Name",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`))
            if err != nil {
                panic(err)
            }

            params := handlers.NewNotifyUserParams(request)
            params.ParseRequestBody()

            Expect(params.Kind).To(Equal("test_email"))
            Expect(params.KindDescription).To(Equal("Descriptive Email Name"))
            Expect(params.SourceDescription).To(Equal("Descriptive Component Name"))
            Expect(params.Subject).To(Equal("Summary of contents"))
            Expect(params.Text).To(Equal("Contents of the email message"))
        })

        It("does not blow up if the request body is empty", func() {
            request, err := http.NewRequest("GET", "/users", strings.NewReader(""))
            if err != nil {
                panic(err)
            }

            params := handlers.NewNotifyUserParams(request)
            Expect(func() {
                params.ParseRequestBody()
            }).NotTo(Panic())
        })
    })

    Describe("ValidateRequestBody", func() {
        It("validates the required parameters in the request body", func() {
            request, err := http.NewRequest("GET", "/users", strings.NewReader(`{
                "kind": "test_email",
                "kind_description": "Descriptive Email Name",
                "source_description": "Descriptive Component Name",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`))
            if err != nil {
                panic(err)
            }
            params := handlers.NewNotifyUserParams(request)
            params.ParseRequestBody()

            Expect(params.ValidateRequestBody()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))

            params.Kind = ""

            Expect(params.ValidateRequestBody()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(1))
            Expect(params.Errors).To(ContainElement(`"kind" is a required field`))

            params.Text = ""

            Expect(params.ValidateRequestBody()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(2))
            Expect(params.Errors).To(ContainElement(`"kind" is a required field`))
            Expect(params.Errors).To(ContainElement(`"text" is a required field`))

            params.Kind = "something"
            params.Text = "banana"

            Expect(params.ValidateRequestBody()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))
        })
    })

    Describe("ParseAuthorizationToken", func() {
        var token string

        BeforeEach(func() {
            tokenHeader := map[string]interface{}{
                "alg": "RS256",
            }
            tokenClaims := map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "my-client",
                "cid":       "my-client",
                "exp":       3404281214,
            }
            token = BuildToken(tokenHeader, tokenClaims)
        })

        It("parses the Authorization header, storing the client_id value", func() {
            request, err := http.NewRequest("GET", "/users", nil)
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)

            params := handlers.NewNotifyUserParams(request)
            params.ParseAuthorizationToken()

            Expect(params.ClientID).To(Equal("my-client"))
        })
    })

    Describe("ValidateAuthorizationToken", func() {
        var request *http.Request
        var err error

        BeforeEach(func() {
            request, err = http.NewRequest("GET", "/users", nil)
            if err != nil {
                panic(err)
            }
        })

        It("validates the presence of an auth token", func() {
            params := handlers.NewNotifyUserParams(request)

            Expect(params.ValidateAuthorizationToken()).To(BeFalse())
            Expect(params.Errors).To(ContainElement("Authorization header is invalid: missing"))
        })

        It("validates the fields of the auth token", func() {
            tokenHeader := map[string]interface{}{
                "alg": "RS256",
            }
            tokenClaims := map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "my-client",
                "cid":       "my-client",
                "exp":       3404281214,
            }
            token := BuildToken(tokenHeader, tokenClaims)

            request.Header.Set("Authorization", "Bearer "+token)

            params := handlers.NewNotifyUserParams(request)
            params.ParseAuthorizationToken()

            Expect(params.ValidateAuthorizationToken()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))

            tokenClaims = map[string]interface{}{
                "jti": "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "cid": "my-client",
                "exp": 3404281214,
            }
            token = BuildToken(tokenHeader, tokenClaims)

            request.Header.Set("Authorization", "Bearer "+token)
            params = handlers.NewNotifyUserParams(request)
            params.ParseAuthorizationToken()

            Expect(params.ValidateAuthorizationToken()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(1))
            Expect(params.Errors).To(ContainElement(`Authorization header is invalid: missing "client_id" field`))

            tokenClaims = map[string]interface{}{
                "jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
                "client_id": "my-client",
                "cid":       "my-client",
                "exp":       1404281214,
            }
            token = BuildToken(tokenHeader, tokenClaims)

            request, err = http.NewRequest("GET", "/users", nil)
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)

            params = handlers.NewNotifyUserParams(request)
            params.ParseAuthorizationToken()

            Expect(params.ValidateAuthorizationToken()).To(BeFalse())
            Expect(len(params.Errors)).To(Equal(1))
            Expect(params.Errors).To(ContainElement(`Authorization header is invalid: expired`))
        })
    })

    Describe("ConfirmPermissions", func() {
        It("validates the scopes in the auth token", func() {
            request, err := http.NewRequest("GET", "/users", nil)
            if err != nil {
                panic(err)
            }

            tokenHeader := map[string]interface{}{
                "alg": "RS256",
            }
            tokenClaims := map[string]interface{}{
                "client_id": "my-client",
                "cid":       "my-client",
                "scope":     []string{},
            }
            token := BuildToken(tokenHeader, tokenClaims)

            request.Header.Set("Authorization", "Bearer "+token)

            params := handlers.NewNotifyUserParams(request)
            Expect(params.ConfirmPermissions()).To(BeFalse())
            Expect(params.Errors).To(ContainElement("You are not authorized to perform the requested action"))

            tokenClaims = map[string]interface{}{
                "client_id": "my-client",
                "cid":       "my-client",
                "scope":     []string{"notifications.write"},
            }
            token = BuildToken(tokenHeader, tokenClaims)
            request.Header.Set("Authorization", "Bearer "+token)

            params = handlers.NewNotifyUserParams(request)
            Expect(params.ConfirmPermissions()).To(BeTrue())
            Expect(len(params.Errors)).To(Equal(0))
        })
    })

    Describe("ParseRequestPath", func() {
        It("reads the user_id out of the request path", func() {
            request, err := http.NewRequest("GET", "/users/my-user-id", nil)
            if err != nil {
                panic(err)
            }

            params := handlers.NewNotifyUserParams(request)
            params.ParseRequestPath()

            Expect(params.UserID).To(Equal("my-user-id"))
        })
    })

    Describe("ParseEnvironmentVariables", func() {
        var sender string

        BeforeEach(func() {
            sender = os.Getenv("SENDER")
        })

        It("reads the SENDER env var into the params object", func() {
            os.Setenv("SENDER", "my-user@example.com")

            request, err := http.NewRequest("GET", "/users/my-user-id", nil)
            if err != nil {
                panic(err)
            }

            params := handlers.NewNotifyUserParams(request)
            params.ParseEnvironmentVariables()

            Expect(params.From).To(Equal("my-user@example.com"))

            os.Setenv("SENDER", sender)
        })
    })
})
