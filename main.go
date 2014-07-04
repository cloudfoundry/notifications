package main

import (
    "log"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/web"
)

func main() {
    configure()
    server := web.NewServer()
    server.Run()
}

func configure() {
    env := config.NewEnvironment()

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
