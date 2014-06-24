package main

import "github.com/pivotal-cf/cf-notifications/web"

func main() {
    server := web.NewServer()
    server.Run()
}
