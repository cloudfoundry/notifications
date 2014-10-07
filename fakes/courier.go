package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type FakeCourier struct {
    Error             error
    Responses         []postal.Response
    DispatchArguments []interface{}
    TheMailer         *FakeMailer
}

func NewFakeCourier() *FakeCourier {
    return &FakeCourier{
        Responses:         make([]postal.Response, 0),
        DispatchArguments: make([]interface{}, 0),
        TheMailer:         NewFakeMailer(),
    }
}

func (fake *FakeCourier) Mailer() postal.MailerInterface {
    return fake.TheMailer
}

func (fake *FakeCourier) Dispatch(token string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {
    fake.DispatchArguments = []interface{}{token, guid, options}
    return fake.Responses, fake.Error
}
