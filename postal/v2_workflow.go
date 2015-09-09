package postal

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/pivotal-golang/lager"
)

type V2Workflow struct{}

func (w V2Workflow) Deliver(job *gobble.Job, logger lager.Logger) {

}
