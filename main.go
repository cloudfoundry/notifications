package main

import "github.com/cloudfoundry-incubator/notifications/application"

func main() {
	app := application.NewApplication()
	defer app.Crash()

	app.Boot()
}
