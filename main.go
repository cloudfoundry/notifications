package main

import "github.com/cloudfoundry-incubator/notifications/web"

func main() {
    server := web.NewServer()
    server.Run()
}
