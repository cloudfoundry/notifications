package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type FakeMailer struct {
    DeliverArguments map[string]interface{}
    Responses        []postal.Response
}

func NewFakeMailer() *FakeMailer {
    return &FakeMailer{}
}

func (fake *FakeMailer) Deliver(conn models.ConnectionInterface, template postal.Templates, users map[string]uaa.User, options postal.Options, space cf.CloudControllerSpace, org cf.CloudControllerOrganization, client string) []postal.Response {
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
