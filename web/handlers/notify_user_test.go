package handlers_test

import (
    "bytes"
    "log"
    "net/http"
    "net/http/httptest"
    "os"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
    var fakeUAA *httptest.Server

    BeforeEach(func() {
        fakeUAA = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            if req.URL.Path == "/Users/user-123" &&
                req.Method == "GET" &&
                req.Header.Get("Authorization") == "Bearer a-special-token" {

                response := `{
                      "id": "user-123",
                      "meta": {
                        "version": 6,
                        "created": "2014-05-22T22:36:36.941Z",
                        "lastModified": "2014-06-25T23:10:03.845Z"
                      },
                      "userName": "admin",
                      "name": {
                        "familyName": "Admin",
                        "givenName": "Mister"
                      },
                      "emails": [
                        {
                          "value": "fake-user@example.com"
                        }
                      ],
                      "groups": [
                        {
                          "value": "e7f74565-4c7e-44ba-b068-b16072cbf08f",
                          "display": "clients.read",
                          "type": "DIRECT"
                        }
                      ],
                      "approvals": [],
                      "active": true,
                      "verified": false,
                      "schemas": [
                        "urn:scim:schemas:core:1.0"
                      ]
                    }`

                w.WriteHeader(http.StatusOK)
                w.Write([]byte(response))
            } else {
                w.WriteHeader(http.StatusNotFound)
            }

        }))

        os.Setenv("UAA_HOST", fakeUAA.URL)
        os.Setenv("UAA_CLIENT_ID", "notifications")
        os.Setenv("UAA_CLIENT_SECRET", "secret")
    })

    AfterEach(func() {
        fakeUAA.Close()
    })

    It("logs the email address of the recipient", func() {
        buffer := bytes.NewBuffer([]byte{})
        logger := log.New(buffer, "", 0)

        writer := httptest.NewRecorder()
        request, err := http.NewRequest("POST", "/users/user-123", nil)
        if err != nil {
            panic(err)
        }
        request.Header.Set("Authorization", "Bearer a-special-token")

        handler := handlers.NewNotifyUser(logger)

        handler.ServeHTTP(writer, request)

        Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
    })
})
