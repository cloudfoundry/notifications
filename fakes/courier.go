package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type Courier struct {
    Error             error
    Responses         []postal.Response
    DispatchArguments []interface{}
    TheMailer         *Mailer
}

func NewCourier() *Courier {
    return &Courier{
        Responses:         make([]postal.Response, 0),
        DispatchArguments: make([]interface{}, 0),
        TheMailer:         NewMailer(),
    }
}

func (fake *Courier) Mailer() postal.MailerInterface {
    return fake.TheMailer
}

func (fake *Courier) Dispatch(token string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {
    fake.DispatchArguments = []interface{}{token, guid, options}
    return fake.Responses, fake.Error
}
