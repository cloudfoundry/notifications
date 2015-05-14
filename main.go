package main

import "github.com/cloudfoundry-incubator/notifications/application"

func main() {
	mother := application.NewMother()
	app := application.NewApplication(mother)
	defer app.Crash()

	app.Boot()
}
