package uaa_test

import (
    "bufio"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
    var client uaa.Client

    BeforeEach(func() {
        client = uaa.Client{}
    })

    Describe("NewClient", func() {
        It("returns a minimally configured Client instance", func() {
            client = uaa.NewClient("http://uaa.example.com", true)
            Expect(client.Host).To(Equal("http://uaa.example.com"))
            Expect(client.BasicAuthUsername).To(Equal(""))
            Expect(client.BasicAuthPassword).To(Equal(""))
            Expect(client.VerifySSL).To(BeTrue())
        })
    })

    Describe("WithBasicAuthCredentials", func() {
        It("returns a client that has Basic Auth credentials set", func() {
            client.AccessToken = "bad-token"

            client = client.WithBasicAuthCredentials("client-id", "the client secret")
            Expect(client.Host).To(Equal(""))
            Expect(client.BasicAuthUsername).To(Equal("client-id"))
            Expect(client.BasicAuthPassword).To(Equal("the client secret"))
            Expect(client.AccessToken).To(Equal(""))
            Expect(client.VerifySSL).To(BeFalse())
        })
    })

    Describe("WithAuthorizationToken", func() {
        It("returns a client that has an authorization token set", func() {
            client.BasicAuthUsername = "bad-user"
            client.BasicAuthPassword = "bad-password"

            client = client.WithAuthorizationToken("token")
            Expect(client.AccessToken).To(Equal("token"))
            Expect(client.BasicAuthUsername).To(Equal(""))
            Expect(client.BasicAuthPassword).To(Equal(""))
        })
    })

    Describe("TLSConfig", func() {
        Context("when VerifySSL option is true", func() {
            It("uses a TLS config that verifies SSL", func() {
                client = uaa.NewClient("", true)
                Expect(client.TLSConfig().InsecureSkipVerify).To(BeFalse())
            })
        })

        Context("when VerifySSL option is false", func() {
            It("uses a TLS config that does not verify SSL", func() {
                client = uaa.NewClient("", false)
                Expect(client.TLSConfig().InsecureSkipVerify).To(BeTrue())
            })
        })
    })

    Describe("MakeRequest", func() {
        var server *httptest.Server
        var headers http.Header

        BeforeEach(func() {
            server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                w.WriteHeader(222)
                buffer := bufio.NewReader(req.Body)
                body, _, err := buffer.ReadLine()
                if err != nil {
                    panic(err)
                }
                response := fmt.Sprintf("%s %s %s", req.Method, req.URL.Path, body)
                headers = req.Header
                w.Write([]byte(response))
            }))

        })

        Context("with basic auth headers", func() {
            It("makes an HTTP request with the given URL, HTTP method, and request body", func() {
                defer server.Close()

                client = uaa.NewClient(server.URL, false).WithBasicAuthCredentials("my-user", "my-pass")

                requestBody := strings.NewReader("key=value")
                code, body, err := client.MakeRequest("GET", "/something", requestBody)
                if err != nil {
                    panic(err)
                }

                Expect(code).To(Equal(222))

                bodyText := string(body)
                Expect(bodyText).To(ContainSubstring("GET"))
                Expect(bodyText).To(ContainSubstring("/something"))
                Expect(bodyText).To(ContainSubstring("key=value"))
                Expect(headers["Content-Type"]).To(ContainElement("application/x-www-form-urlencoded"))
                Expect(strings.Join(headers["Authorization"], " ")).To(ContainSubstring("Basic bXktdXNlcjpteS1wYXNz"))
            })
        })

        Context("with oaut access token", func() {
            It("makes an HTTP request with the given URL, HTTP method, and request body", func() {
                defer server.Close()

                client = uaa.NewClient(server.URL, false)
                client = client.WithBasicAuthCredentials("my-user", "my-pass")
                client = client.WithAuthorizationToken("my-special-token")

                requestBody := strings.NewReader("key=value")
                code, body, err := client.MakeRequest("GET", "/something", requestBody)
                if err != nil {
                    panic(err)
                }

                Expect(code).To(Equal(222))

                bodyText := string(body)
                Expect(bodyText).To(ContainSubstring("GET"))
                Expect(bodyText).To(ContainSubstring("/something"))
                Expect(bodyText).To(ContainSubstring("key=value"))
                Expect(headers["Content-Type"]).To(ContainElement("application/x-www-form-urlencoded"))
                Expect(strings.Join(headers["Authorization"], " ")).To(ContainSubstring("Bearer my-special-token"))
            })
        })
    })
})
