package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
)

type Validator struct {
	ValidateCall struct {
		Receives struct {
			Params *notify.NotifyParams
		}
		Returns struct {
			Valid bool
		}
		ErrorsToApply []string
	}
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(params *notify.NotifyParams) bool {
	v.ValidateCall.Receives.Params = params
	params.Errors = append(params.Errors, v.ValidateCall.ErrorsToApply...)

	return v.ValidateCall.Returns.Valid
}
