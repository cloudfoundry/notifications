package handlers_test

import (
    "encoding/base64"
    "strings"
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Suite")
}

type Envelope struct {
    AuthUser string
    AuthPass string
    From     string
    To       string
    Data     []string
}

func (envelope *Envelope) Respond(request string) (string, bool) {
    switch {
    case strings.Contains(request, "EHLO"):
        return "250-localhost Hello\n250-SIZE 52428800\n250-PIPELINING\n250-AUTH PLAIN LOGIN\n250 HELP", false
    case strings.Contains(request, "AUTH"):
        auth := strings.TrimPrefix(request, "AUTH PLAIN ")
        decoded, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(auth))
        parts := strings.SplitN(string(decoded), "\x00", 3)
        envelope.AuthUser = parts[1]
        envelope.AuthPass = parts[2]
        return "235 OK, Go ahead", false
    case strings.Contains(request, "MAIL FROM"):
        from := strings.TrimPrefix(request, "MAIL FROM:")
        envelope.From = strings.TrimSpace(from)
        return "250 OK", false
    case strings.Contains(request, "RCPT TO"):
        to := strings.TrimPrefix(request, "RCPT TO:")
        envelope.To = strings.TrimSpace(to)
        return "250 OK", false
    case strings.Contains(request, "DATA"):
        return "354 Go ahead", false
    case strings.TrimSpace(request) == ".":
        if strings.Contains(envelope.To, "bounce") {
            return "450 Mailbox unavailable", false
        }
        return "250 Written safely to disk", false
    case strings.Contains(request, "QUIT"):
        return "221 localhost saying goodbye", true
    default:
        envelope.Data = append(envelope.Data, strings.TrimSpace(request))
        return "", false
    }
}
