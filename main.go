package main

import (
    "errors"
    "log"
    "net"
    "time"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

func main() {
    defer crash()

    env := config.NewEnvironment()
    configure(env)

    if !env.TestMode {
        confirmSMTPConfiguration(env)
    }
    retrieveUAAPublicKey()
    migrate()

    server := web.NewServer()
    server.Run()
}

func configure(env config.Environment) {
    log.Println("Booting with configuration:")
    log.Printf("\tCCHost          -> %+v", env.CCHost)
    log.Printf("\tDatabaseURL     -> %+v", env.DatabaseURL)
    log.Printf("\tSMTPHost        -> %+v", env.SMTPHost)
    log.Printf("\tSMTPPass        -> %+v", env.SMTPPass)
    log.Printf("\tSMTPPort        -> %+v", env.SMTPPort)
    log.Printf("\tSMTPTLS         -> %+v", env.SMTPTLS)
    log.Printf("\tSMTPUser        -> %+v", env.SMTPUser)
    log.Printf("\tSender          -> %+v", env.Sender)
    log.Printf("\tTest mode       -> %+v", env.TestMode)
    log.Printf("\tUAAClientID     -> %+v", env.UAAClientID)
    log.Printf("\tUAAClientSecret -> %+v", env.UAAClientSecret)
    log.Printf("\tUAAHost         -> %+v", env.UAAHost)
    log.Printf("\tVerifySSL       -> %+v", env.VerifySSL)
}

func confirmSMTPConfiguration(env config.Environment) {
    mailClient, err := mail.NewClient(env.SMTPUser, env.SMTPPass, net.JoinHostPort(env.SMTPHost, env.SMTPPort))
    if err != nil {
        panic(err)
    }

    err = mailClient.Connect()
    if err != nil {
        panic(err)
    }

    err = mailClient.Hello()
    if err != nil {
        panic(err)
    }

    startTLSSupported, _ := mailClient.Extension("STARTTLS")

    mailClient.Quit()

    if !startTLSSupported && env.SMTPTLS {
        panic(errors.New(`SMTP TLS configuration mismatch: Configured to use TLS over SMTP, but the mail server does not support the "STARTTLS" extension.`))
    }

    if startTLSSupported && !env.SMTPTLS {
        panic(errors.New(`SMTP TLS configuration mismatch: Not configured to use TLS over SMTP, but the mail server does support the "STARTTLS" extension.`))
    }
}

func retrieveUAAPublicKey() {
    env := config.NewEnvironment()
    auth := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
    auth.VerifySSL = env.VerifySSL

    key, err := uaa.GetTokenKey(auth)
    if err != nil {
        panic(err)
    }

    config.UAAPublicKey = key
    log.Printf("UAA Public Key: %s", config.UAAPublicKey)
}

func migrate() {
    models.Database()
}

// This is a hack to get the logs to output to the loggregator before the process exits
func crash() {
    err := recover()
    if err != nil {
        time.Sleep(5 * time.Second)
        panic(err)
    }
}
