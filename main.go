package main

import "github.com/cloudfoundry-incubator/notifications/application"

func main() {
    app := application.NewApplication()
    defer app.Crash()

    app.PrintConfiguration()
    app.ConfigureSMTP()
    app.RetrieveUAAPublicKey()
    app.Migrate()
    app.EnableDBLogging()
    app.UnlockJobs()
    app.StartWorkers()
    app.StartServer()
}
