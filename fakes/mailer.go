package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/postal/strategies"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type Mailer struct {
    DeliverArguments map[string]interface{}
    Responses        []strategies.Response
}

func NewMailer() *Mailer {
    return &Mailer{}
}

func (fake *Mailer) Deliver(conn models.ConnectionInterface, template postal.Templates, users map[string]uaa.User, options postal.Options, space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client string) []strategies.Response {
    fake.DeliverArguments = map[string]interface{}{
        "connection": conn,
        "template":   template,
        "users":      users,
        "options":    options,
        "space":      space,
        "org":        org,
        "client":     client,
    }

    return fake.Responses
}
