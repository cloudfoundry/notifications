package config

import (
    "errors"
    "fmt"
    "os"
    "strconv"
)

var UAAPublicKey string

type Environment struct {
    UAAHost         string
    UAAClientID     string
    UAAClientSecret string
    SMTPUser        string
    SMTPPass        string
    SMTPHost        string
    SMTPPort        string
    SMTPTLS         bool
    Sender          string
    CCHost          string
}

func NewEnvironment() Environment {
    return Environment{
        UAAHost:         loadOrPanic("UAA_HOST"),
        UAAClientID:     loadOrPanic("UAA_CLIENT_ID"),
        UAAClientSecret: loadOrPanic("UAA_CLIENT_SECRET"),
        SMTPUser:        loadOrPanic("SMTP_USER"),
        SMTPPass:        loadOrPanic("SMTP_PASS"),
        SMTPHost:        loadOrPanic("SMTP_HOST"),
        SMTPPort:        loadOrPanic("SMTP_PORT"),
        SMTPTLS:         loadBool("SMTP_TLS"),
        Sender:          loadOrPanic("SENDER"),
        CCHost:          loadOrPanic("CC_HOST"),
    }
}

func loadOrPanic(name string) string {
    value := os.Getenv(name)
    if value == "" {
        panic(errors.New(fmt.Sprintf("Could not find required %s environment variable", name)))
    }
    return value
}

func loadBool(name string) bool {
    if os.Getenv(name) == "" {
        return false
    }
    value, err := strconv.ParseBool(os.Getenv(name))
    if err != nil {
        panic(errors.New(fmt.Sprintf("Could not parse %s environment variable into boolean", name)))
    }
    return value
}
