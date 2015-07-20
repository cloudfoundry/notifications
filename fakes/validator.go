package fakes

import "github.com/cloudfoundry-incubator/notifications/web/v1/notify"

type Validator struct {
	ValidateErrors []string
}

func (fake *Validator) Validate(n *notify.NotifyParams) bool {
	n.Errors = append(n.Errors, fake.ValidateErrors...)
	return len(fake.ValidateErrors) == 0
}
