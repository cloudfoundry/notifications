package application

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/ryanmoran/viron"
)

type Environment struct {
	CCHost                             string `env:"CC_HOST" env-required:"true"`
	CORSOrigin                         string `env:"CORS_ORIGIN" env-default:"*"`
	DBLoggingEnabled                   bool   `env:"DB_LOGGING_ENABLED"`
	DBMaxOpenConns                     int    `env:"DB_MAX_OPEN_CONNS"`
	DatabaseURL                        string `env:"DATABASE_URL" env-required:"true"`
	DefaultUAAScopesList               string `env:"DEFAULT_UAA_SCOPES"`
	Domain                             string `env:"DOMAIN" env-required:"true"`
	EncryptionKey                      []byte `env:"ENCRYPTION_KEY" env-required:"true"`
	GobbleWaitMaxDuration              int    `env:"GOBBLE_WAIT_MAX_DURATION" env-default:"5000"`
	GobbleMaxQueueLength               int    `env:"GOBBLE_MAX_QUEUE_LENGTH" env-default:"5000"`
	MaxRetries                         int    `env:"MAX_RETRIES" env-default:"5"`
	Port                               int    `env:"PORT" env-default:"3000"`
	RootPath                           string `env:"ROOT_PATH"`
	SMTPAuthMechanism                  string `env:"SMTP_AUTH_MECHANISM" env-required:"true"`
	SMTPCRAMMD5Secret                  string `env:"SMTP_CRAMMD5_SECRET"`
	SMTPHost                           string `env:"SMTP_HOST" env-required:"true"`
	SMTPLoggingEnabled                 bool   `env:"SMTP_LOGGING_ENABLED" env-default:"false"`
	SMTPPass                           string `env:"SMTP_PASS"`
	SMTPPort                           string `env:"SMTP_PORT" env-required:"true"`
	SMTPTLS                            bool   `env:"SMTP_TLS" env-default:"true"`
	SMTPUser                           string `env:"SMTP_USER"`
	Sender                             string `env:"SENDER" env-required:"true"`
	TestMode                           bool   `env:"TEST_MODE" env-default:"false"`
	UAAClientID                        string `env:"UAA_CLIENT_ID" env-required:"true"`
	UAAClientSecret                    string `env:"UAA_CLIENT_SECRET" env-required:"true"`
	UAAHost                            string `env:"UAA_HOST" env-required:"true"`
	UAAKeyRefreshInterval              int    `env:"UAA_KEY_REFRESH_INTREVAL" env-default:"60000"`
	VerifySSL                          bool   `env:"VERIFY_SSL" env-default:"true"`
	DatabaseCACertFile                 string `env:"DATABASE_CA_CERT_FILE"`
	DatabaseCommonName                 string `env:"DATABASE_COMMON_NAME"`
	DatabaseEnableIdentityVerification bool   `env:"DATABASE_ENABLE_IDENTITY_VERIFICATION" env-default:"true"`

	VCAPApplication struct {
		InstanceIndex int `json:"instance_index"`
	} `env:"VCAP_APPLICATION" env-required:"true"`

	ModelMigrationsPath  string
	GobbleMigrationsPath string
	DefaultUAAScopes     []string
}

type EnvironmentError struct {
	Err error
}

func (e EnvironmentError) Error() string {
	return e.Err.Error() + " (Please see https://github.com/cloudfoundry/notifications-release to find a packaged version of notifications and see the required configuration)"
}

func NewEnvironment() (Environment, error) {
	env := Environment{}
	err := viron.Parse(&env)
	if err != nil {
		return env, EnvironmentError{err}
	}

	err = env.parseDatabaseURL()
	if err != nil {
		return env, EnvironmentError{err}
	}

	env.expandRoot()

	err = env.validateSMTPAuthMechanism()
	if err != nil {
		return env, EnvironmentError{err}
	}

	env.inferMigrationsDirs()
	env.parseDefaultUAAScopes()

	return env, nil
}

func (env *Environment) parseDefaultUAAScopes() {
	env.DefaultUAAScopes = strings.Split(env.DefaultUAAScopesList, ",")
}

func (env *Environment) expandRoot() {
	env.RootPath = os.ExpandEnv(env.RootPath)
}

func (env *Environment) inferMigrationsDirs() {
	env.ModelMigrationsPath = path.Join(env.RootPath, "db", "migrations")
	env.GobbleMigrationsPath = path.Join(env.RootPath, "gobble", "migrations")
}

func (env *Environment) parseDatabaseURL() error {
	databaseURL := env.DatabaseURL
	databaseURL = strings.TrimPrefix(databaseURL, "http://")
	databaseURL = strings.TrimPrefix(databaseURL, "https://")
	databaseURL = strings.TrimPrefix(databaseURL, "tcp://")
	databaseURL = strings.TrimPrefix(databaseURL, "mysql://")
	databaseURL = strings.TrimPrefix(databaseURL, "mysql2://")
	parsedURL, err := url.Parse("tcp://" + databaseURL)
	if err != nil {
		return fmt.Errorf("Could not parse DATABASE_URL %q, it does not fit format %q", env.DatabaseURL, "tcp://user:pass@host/dname")
	}

	password, _ := parsedURL.User.Password()
	env.DatabaseURL = fmt.Sprintf("%s:%s@%s(%s)%s?parseTime=true", parsedURL.User.Username(), password, parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
	return nil
}

func (env *Environment) validateSMTPAuthMechanism() error {
	for _, mechanism := range mail.SMTPAuthMechanisms {
		if mechanism == env.SMTPAuthMechanism {
			return nil
		}
	}

	return fmt.Errorf("Could not parse SMTP_AUTH_MECHANISM %q, it is not one of the allowed values: %+v", env.SMTPAuthMechanism, mail.SMTPAuthMechanisms)
}
