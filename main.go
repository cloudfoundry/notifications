package main

import "github.com/cloudfoundry-incubator/notifications/application"

func main() {
	app := application.NewApplication()
	logger := application.BootLogger()
	defer app.Crash(logger)

	app.Boot(logger)
}
