package config

import (
    "errors"
    "fmt"
    "net/url"
    "os"
    "strconv"
    "strings"
)

var UAAPublicKey string

type Environment struct {
    CCHost          string
    DatabaseURL     string
    Port            string
    RootPath        string
    SMTPHost        string
    SMTPPass        string
    SMTPPort        string
    SMTPTLS         bool
    SMTPUser        string
    Sender          string
    UAAClientID     string
    UAAClientSecret string
    UAAHost         string
    VerifySSL       bool
    TestMode        bool
}

func NewEnvironment() Environment {
    return Environment{
        CCHost:          loadOrPanic("CC_HOST"),
        DatabaseURL:     loadDatabaseURL("DATABASE_URL"),
        Port:            loadPort(),
        RootPath:        loadOrPanic("ROOT_PATH"),
        SMTPHost:        loadOrPanic("SMTP_HOST"),
        SMTPPass:        loadOrPanic("SMTP_PASS"),
        SMTPPort:        loadOrPanic("SMTP_PORT"),
        SMTPTLS:         loadBool("SMTP_TLS", true),
        SMTPUser:        loadOrPanic("SMTP_USER"),
        Sender:          loadOrPanic("SENDER"),
        UAAClientID:     loadOrPanic("UAA_CLIENT_ID"),
        UAAClientSecret: loadOrPanic("UAA_CLIENT_SECRET"),
        UAAHost:         loadOrPanic("UAA_HOST"),
        VerifySSL:       loadBool("VERIFY_SSL", true),
        TestMode:        loadBool("TEST_MODE", false),
    }
}

func loadOrPanic(name string) string {
    value := os.Getenv(name)
    if value == "" {
        panic(errors.New(fmt.Sprintf("Could not find required %s environment variable", name)))
    }
    return value
}

func loadDatabaseURL(name string) string {
    databaseURL := loadOrPanic(name)
    databaseURL = strings.TrimPrefix(databaseURL, "http://")
    databaseURL = strings.TrimPrefix(databaseURL, "https://")
    databaseURL = strings.TrimPrefix(databaseURL, "tcp://")
    parsedURL, err := url.Parse("tcp://" + databaseURL)
    if err != nil {
        panic(errors.New(fmt.Sprintf("Could not parse DATABASE_URL %q, it does not fit format %q", loadOrPanic(name), "tcp://user:pass@host/dname")))
    }

    password, _ := parsedURL.User.Password()
    return fmt.Sprintf("%s:%s@%s(%s)%s?parseTime=true", parsedURL.User.Username(), password, parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
}

func loadPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        return "3000"
    }
    return port
}

func loadBool(name string, defaultValue bool) bool {
    value, err := strconv.ParseBool(os.Getenv(name))
    if err != nil {
        return defaultValue
    }

    return value
}
