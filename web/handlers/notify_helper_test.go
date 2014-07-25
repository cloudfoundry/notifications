package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyHelper", func() {
    var helper handlers.NotifyHelper
    var fakeCC *FakeCloudController
    var logger *log.Logger
    var request *http.Request
    var fakeUAA FakeUAAClient
    var mailClient FakeMailClient
    var writer *httptest.ResponseRecorder
    var token string
    var buffer *bytes.Buffer

    BeforeEach(func() {

        tokenHeader := map[string]interface{}{
            "alg": "FAST",
        }

        tokenClaims := map[string]interface{}{
            "client_id": "mister-client",
            "exp":       3404281214,
            "scope":     []string{"notifications.write"},
        }

        token = BuildToken(tokenHeader, tokenClaims)

        writer = httptest.NewRecorder()

        fakeCC = NewFakeCloudController()
        fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
            cf.CloudControllerUser{Guid: "user-123"},
            cf.CloudControllerUser{Guid: "user-456"},
        }

        fakeCC.Spaces = map[string]cf.CloudControllerSpace{
            "space-001": cf.CloudControllerSpace{
                Name:             "production",
                Guid:             "space-001",
                OrganizationGuid: "org-001",
            },
        }

        fakeCC.Orgs = map[string]cf.CloudControllerOrganization{
            "org-001": cf.CloudControllerOrganization{
                Name: "pivotaltracker",
            },
        }

        fakeUAA = FakeUAAClient{
            ClientToken: uaa.Token{
                Access: token,
            },
            UsersByID: map[string]uaa.User{
                "user-123": uaa.User{
                    ID:     "user-123",
                    Emails: []string{"user-123@example.com"},
                },
                "user-456": uaa.User{
                    ID:     "user-456",
                    Emails: []string{"user-456@example.com"},
                },
            },
        }

        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)

        mailClient = FakeMailClient{}

        helper = handlers.NewNotifyHelper(fakeCC, logger, &fakeUAA, FakeGuidGenerator, &mailClient)
    })

    Describe("NofifyServeHTTP", func() {
        var loadCCUsers func(string, string) ([]cf.CloudControllerUser, error)

        BeforeEach(func() {
            loadCCUsers = func(userGUID, accessToken string) ([]cf.CloudControllerUser, error) {
                return []cf.CloudControllerUser{
                    cf.CloudControllerUser{Guid: userGUID},
                }, nil
            }
        })

        Context("when the request is invalid", func() {
            BeforeEach(func() {
                requestBody, err := json.Marshal(map[string]string{
                    "kind_description":   "Password reminder",
                    "source_description": "Login system",
                    "subject":            "Reset your password",
                })

                if err != nil {
                    panic(err)
                }

                request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
                if err != nil {
                    panic(err)
                }
            })

            It("returns an error message", func() {
                helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                parsed := map[string][]string{}
                err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                if err != nil {
                    panic(err)
                }

                Expect(parsed["errors"]).To(ContainElement(`"kind" is a required field`))
                Expect(parsed["errors"]).To(ContainElement(`"text" or "html" fields must be supplied`))
            })

            Context("when the request body is missing", func() {
                BeforeEach(func() {
                    request.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
                })

                It("returns an error message", func() {
                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    parsed := map[string][]string{}
                    err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parsed["errors"]).To(ContainElement(`"kind" is a required field`))
                    Expect(parsed["errors"]).To(ContainElement(`"text" or "html" fields must be supplied`))
                })
            })

        })

        Context("when the request is valid", func() {
            BeforeEach(func() {
                requestBody, err := json.Marshal(map[string]string{
                    "kind":               "forgot_password",
                    "kind_description":   "Password reminder",
                    "source_description": "Login system",
                    "text":               "Please reset your password by clicking on this link...",
                    "html":               "<p>Please reset your password by clicking on this link...</p>",
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

            Context("when the SMTP server fails to deliver the mail", func() {
                It("returns a status indicating that delivery failed", func() {
                    mailClient.errorOnSend = true
                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    Expect(writer.Code).To(Equal(http.StatusOK))
                    parsed := []map[string]string{}
                    err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parsed[0]["status"]).To(Equal("failed"))
                })
            })

            Context("when the SMTP server cannot be reached", func() {
                It("returns a status indicating that the server is unavailable", func() {
                    mailClient.errorOnConnect = true
                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    Expect(writer.Code).To(Equal(http.StatusOK))
                    parsed := []map[string]string{}
                    err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parsed[0]["status"]).To(Equal("unavailable"))
                })
            })

            Context("when UAA cannot be reached", func() {
                It("returns a 502 status code", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))
                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    Expect(writer.Code).To(Equal(http.StatusBadGateway))
                    Expect(writer.Body.String()).To(ContainSubstring("{\"errors\":[\"UAA is unavailable\"]}"))
                })
            })

            Context("when UAA cannot find the user", func() {
                It("returns that the user in the response with status notfound", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("User f3b51aac-866e-4b7a-948c-de31beefc475d does not exist"))
                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    Expect(writer.Code).To(Equal(http.StatusOK))

                    response := []map[string]string{}
                    err := json.Unmarshal(writer.Body.Bytes(), &response)
                    if err != nil {
                        panic(err)
                    }
                    Expect(response[0]["status"]).To(Equal("notfound"))
                    Expect(response[0]["recipient"]).To(Equal("user-123"))
                })
            })

            Context("when the UAA user has no email", func() {
                It("returns the user in the response with the status noaddress", func() {
                    fakeUAA.UsersByID["user-123"] = uaa.User{
                        ID:     "user-123",
                        Emails: []string{},
                    }

                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    response := []map[string]string{}
                    err := json.Unmarshal(writer.Body.Bytes(), &response)
                    if err != nil {
                        panic(err)
                    }

                    Expect(writer.Code).To(Equal(http.StatusOK))
                    Expect(response[0]["status"]).To(Equal("noaddress"))
                })
            })

            Context("when UAA causes some unknown error", func() {
                It("returns a 502 status code", func() {
                    fakeUAA.ErrorForUserByID = errors.New("Boom!")
                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    Expect(writer.Code).To(Equal(http.StatusBadGateway))
                    Expect(writer.Body.String()).To(ContainSubstring("{\"errors\":[\"UAA Unknown Error: Boom!\"]}"))
                })
            })

            Context("When load Users returns multiple users", func() {
                BeforeEach(func() {
                    loadCCUsers = func(userGUID, accessToken string) ([]cf.CloudControllerUser, error) {
                        return []cf.CloudControllerUser{
                            cf.CloudControllerUser{Guid: "user-123"},
                            cf.CloudControllerUser{Guid: "user-456"},
                        }, nil
                    }
                })

                It("logs the UUIDs of all recipients", func() {
                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, true)

                    lines := strings.Split(buffer.String(), "\n")

                    Expect(lines).To(ContainElement("CloudController user guid: user-123"))
                    Expect(lines).To(ContainElement("CloudController user guid: user-456"))
                })

                It("returns necessary info in the response for the sent mail", func() {
                    helper = handlers.NewNotifyHelper(fakeCC, logger, &fakeUAA, func() (*uuid.UUID, error) {
                        guid, err := uuid.NewV4()
                        if err != nil {
                            panic(err)
                        }
                        return guid, nil
                    }, &mailClient)

                    helper.NotifyServeHTTP(writer, request, "user-123", loadCCUsers, false)

                    Expect(writer.Code).To(Equal(http.StatusOK))
                    parsed := []map[string]string{}
                    err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                    if err != nil {
                        panic(err)
                    }

                    Expect(parsed[0]["status"]).To(Equal("delivered"))
                    Expect(parsed[0]["recipient"]).To(Equal("user-123"))
                    Expect(parsed[0]["notification_id"]).NotTo(Equal(""))

                    Expect(parsed[1]["status"]).To(Equal("delivered"))
                    Expect(parsed[1]["recipient"]).To(Equal("user-456"))
                    Expect(parsed[1]["notification_id"]).NotTo(Equal(parsed[0]["notification_id"]))
                })
            })

            Context("when sending emails to a space", func() {
                BeforeEach(func() {
                    loadCCUsers = func(userGUID, accessToken string) ([]cf.CloudControllerUser, error) {
                        return []cf.CloudControllerUser{
                            cf.CloudControllerUser{Guid: "user-123"},
                            cf.CloudControllerUser{Guid: "user-456"},
                        }, nil
                    }
                    requestBody, err := json.Marshal(map[string]string{
                        "kind":               "forgot_password",
                        "kind_description":   "Password reminder",
                        "source_description": "Login system",
                        "text":               "Please reset your password by clicking on this link...",
                        "html":               "<p>Please reset your password by clicking on this link...</p>",
                    })
                    if err != nil {
                        panic(err)
                    }

                    request, err = http.NewRequest("POST", "/space/space-001", bytes.NewReader(requestBody))
                    if err != nil {
                        panic(err)
                    }
                    request.Header.Set("Authorization", "Bearer "+token)
                })
            })
        })
    })

    Describe("LoadUaaUser", func() {
        Context("UAA returns a user", func() {
            It("returns the uaa.User", func() {
                user, err := helper.LoadUaaUser("user-123", fakeUAA)
                if err != nil {
                    panic(err)
                }

                Expect(user.ID).To(Equal("user-123"))
                Expect(user.Emails[0]).To(Equal("user-123@example.com"))
            })
        })

        Describe("UAA Error Responses", func() {
            Context("when UAA cannot be reached", func() {
                It("returns a UAADownError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAADownError{}))
                })
            })

            Context("when UAA cannot find the user", func() {
                It("returns a UAAUserNotFoundError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("User f3b51aac-866e-4b7a-948c-de31beefc475d does not exist"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAAUserNotFoundError{}))
                })
            })

            Context("when UAA returns an unknown UAA 404 error", func() {
                It("returns a UAAGenericError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(404, []byte("Weird message we haven't seen"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAAGenericError{}))
                })
            })

            Context("when UAA returns an failure code that is not 404", func() {
                It("returns a UAADownError", func() {
                    fakeUAA.ErrorForUserByID = uaa.NewFailure(500, []byte("Doesn't matter"))

                    _, err := helper.LoadUaaUser("user-123", fakeUAA)

                    Expect(err).To(BeAssignableToTypeOf(handlers.UAADownError{}))
                })
            })
        })
    })
})
