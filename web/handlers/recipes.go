package handlers

import (
    "encoding/json"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    EmailFieldName      = "email"
    RecipientsFieldName = "recipient"
)

type MailRecipeInterface interface {
    DeliverMail(clientID string, guid postal.TypedGUID,
        options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error)
    Trim([]byte) []byte
}

type EmailRecipe struct {
    courier postal.CourierInterface
}

func NewEmailRecipe(courier postal.CourierInterface) EmailRecipe {
    return EmailRecipe{
        courier: courier,
    }
}

func (recipe EmailRecipe) DeliverMail(clientID string, guid postal.TypedGUID,
    options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {

    users := map[string]uaa.User{"no-guid-yet": uaa.User{Emails: []string{options.To}}}
    space := ""
    organization := ""
    template := postal.Templates{
        Subject: "",
        Text:    "",
        HTML:    "email template",
    }

    return recipe.courier.Mailer().Deliver(conn, template, users, options, space, organization, clientID), nil
}

func (recipe EmailRecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, RecipientsFieldName)
}

type UAARecipe struct {
    courier postal.CourierInterface
}

func NewUAARecipe(courier postal.CourierInterface) UAARecipe {
    return UAARecipe{courier: courier}
}

func (recipe UAARecipe) DeliverMail(clientID string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {
    return recipe.courier.Dispatch(clientID, guid, options, conn)
}

func (recipe UAARecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, EmailFieldName)
}

type Trimmer struct{}

func (t Trimmer) TrimFields(responses []byte, field string) []byte {
    var results []map[string]string

    err := json.Unmarshal(responses, &results)
    if err != nil {
        panic(err)
    }

    for _, value := range results {
        delete(value, field)
    }

    responses, err = json.Marshal(results)
    if err != nil {
        panic(err)
    }

    return responses
}
