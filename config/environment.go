package config

import (
    "errors"
    "fmt"
    "net/url"
    "strings"

    "github.com/ryanmoran/viron"
)

var UAAPublicKey string

type Environment struct {
    CCHost           string `env:"CC_HOST"            env-required:"true"`
    CORSOrigin       string `env:"CORS_ORIGIN"        env-default:"*"`
    DBLoggingEnabled bool   `env:"DB_LOGGING_ENABLED"`
    DatabaseURL      string `env:"DATABASE_URL"       env-required:"true"`
    Port             string `env:"PORT"               env-default:"3000"`
    RootPath         string `env:"ROOT_PATH"`
    SMTPHost         string `env:"SMTP_HOST"          env-required:"true"`
    SMTPPass         string `env:"SMTP_PASS"`
    SMTPPort         string `env:"SMTP_PORT"          env-required:"true"`
    SMTPTLS          bool   `env:"SMTP_TLS"           env-default:"true"`
    SMTPUser         string `env:"SMTP_USER"`
    Sender           string `env:"SENDER"             env-required:"true"`
    TestMode         bool   `env:"TEST_MODE"          env-default:"false"`
    UAAClientID      string `env:"UAA_CLIENT_ID"      env-required:"true"`
    UAAClientSecret  string `env:"UAA_CLIENT_SECRET"  env-required:"true"`
    UAAHost          string `env:"UAA_HOST"           env-required:"true"`
    VerifySSL        bool   `env:"VERIFY_SSL"         env-default:"true"`
    VCAPApplication  struct {
        InstanceIndex int `json:"instance_index"`
    }   `env:"VCAP_APPLICATION" env-required:"true"`
}

func NewEnvironment() Environment {
    env := Environment{}
    err := viron.Parse(&env)
    if err != nil {
        panic(err)
    }
    env.parseDatabaseURL()

    return env
}

func (env *Environment) parseDatabaseURL() {
    databaseURL := env.DatabaseURL
    databaseURL = strings.TrimPrefix(databaseURL, "http://")
    databaseURL = strings.TrimPrefix(databaseURL, "https://")
    databaseURL = strings.TrimPrefix(databaseURL, "tcp://")
    databaseURL = strings.TrimPrefix(databaseURL, "mysql://")
    databaseURL = strings.TrimPrefix(databaseURL, "mysql2://")
    parsedURL, err := url.Parse("tcp://" + databaseURL)
    if err != nil {
        panic(errors.New(fmt.Sprintf("Could not parse DATABASE_URL %q, it does not fit format %q", env.DatabaseURL, "tcp://user:pass@host/dname")))
    }

    password, _ := parsedURL.User.Password()
    env.DatabaseURL = fmt.Sprintf("%s:%s@%s(%s)%s?parseTime=true", parsedURL.User.Username(), password, parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
}
