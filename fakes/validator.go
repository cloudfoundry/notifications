package fakes

import "github.com/cloudfoundry-incubator/notifications/web/params"

type Validator struct {
	ValidateErrors []string
}

func (fake *Validator) Validate(notify *params.NotifyParams) bool {
	notify.Errors = append(notify.Errors, fake.ValidateErrors...)
	return len(fake.ValidateErrors) == 0
}
