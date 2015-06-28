package fakes

import "github.com/cloudfoundry-incubator/notifications/web/handlers"

type Validator struct {
	ValidateErrors []string
}

func (fake *Validator) Validate(notify *handlers.NotifyParams) bool {
	notify.Errors = append(notify.Errors, fake.ValidateErrors...)
	return len(fake.ValidateErrors) == 0
}
