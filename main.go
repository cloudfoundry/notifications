package main

import (
	"log"

	"github.com/cloudfoundry-incubator/notifications/application"
)

func main() {
	env, err := application.NewEnvironment()
	if err != nil {
		log.Fatalf("CRASHING: %s\n", err)
	}

	mother := application.NewMother(env)
	app := application.New(env, mother)
	defer app.Crash()

	app.Run()
}
