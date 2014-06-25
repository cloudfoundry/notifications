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
    log.Printf("\tUAAHost         -> %s", env.UAAHost)
    log.Printf("\tUAAClientID     -> %s", env.UAAClientID)
    log.Printf("\tUAAClientSecret -> %s", env.UAAClientSecret)
}
