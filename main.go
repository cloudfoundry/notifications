package main

import (
    "errors"
    "log"
    "net"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/web"
)

func main() {
    env := config.NewEnvironment()
    configure(env)
    confirmSMTPConfiguration(env)
    server := web.NewServer()
    server.Run()
}

func configure(env config.Environment) {
    log.Println("Booting with configuration:")
    log.Printf("\tUAAHost         -> %+v", env.UAAHost)
    log.Printf("\tUAAClientID     -> %+v", env.UAAClientID)
    log.Printf("\tUAAClientSecret -> %+v", env.UAAClientSecret)
    log.Printf("\tSMTPUser        -> %+v", env.SMTPUser)
    log.Printf("\tSMTPPass        -> %+v", env.SMTPPass)
    log.Printf("\tSMTPHost        -> %+v", env.SMTPHost)
    log.Printf("\tSMTPPort        -> %+v", env.SMTPPort)
    log.Printf("\tSMTPTLS         -> %+v", env.SMTPTLS)
    log.Printf("\tSender          -> %+v", env.Sender)
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
