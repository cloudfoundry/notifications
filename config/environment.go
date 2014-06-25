package config

import (
    "errors"
    "fmt"
    "os"
)

type Environment struct {
    UAAHost         string
    UAAClientID     string
    UAAClientSecret string
}

func NewEnvironment() Environment {
    return Environment{
        UAAHost:         loadOrPanic("UAA_HOST"),
        UAAClientID:     loadOrPanic("UAA_CLIENT_ID"),
        UAAClientSecret: loadOrPanic("UAA_CLIENT_SECRET"),
    }
}

func loadOrPanic(name string) string {
    value := os.Getenv(name)
    if value == "" {
        panic(errors.New(fmt.Sprintf("Could not find required %s environment variable", name)))
    }
    return value
}
