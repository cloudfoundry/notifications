package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/postal/strategies"
)

type Courier struct {
    Error             error
    Responses         []strategies.Response
    DispatchArguments []interface{}
    TheMailer         *Mailer
}

func NewCourier() *Courier {
    return &Courier{
        Responses:         make([]strategies.Response, 0),
        DispatchArguments: make([]interface{}, 0),
        TheMailer:         NewMailer(),
    }
}

func (fake *Courier) Mailer() strategies.MailerInterface {
    return fake.TheMailer
}

func (fake *Courier) Dispatch(token string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]strategies.Response, error) {
    fake.DispatchArguments = []interface{}{token, guid, options}
    return fake.Responses, fake.Error
}
