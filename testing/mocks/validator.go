package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/web/notify"

type Validator struct {
	ValidateErrors []string
}

func (fake *Validator) Validate(n *notify.NotifyParams) bool {
	n.Errors = append(n.Errors, fake.ValidateErrors...)
	return len(fake.ValidateErrors) == 0
}
