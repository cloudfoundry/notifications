package postal

import (
    "encoding/json"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    EmailFieldName      = "email"
    RecipientsFieldName = "recipient"
)

type MailRecipeInterface interface {
    DeliverMail(clientID string, guid TypedGUID,
        options Options, conn models.ConnectionInterface) ([]Response, error)
    Trim([]byte) []byte
}

type EmailRecipe struct {
    courier        CourierInterface
    templateLoader TemplateLoaderInterface
}

func NewEmailRecipe(courier CourierInterface, templateLoader TemplateLoaderInterface) EmailRecipe {
    return EmailRecipe{
        courier:        courier,
        templateLoader: templateLoader,
    }
}

func (recipe EmailRecipe) DeliverMail(clientID string, guid TypedGUID,
    options Options, conn models.ConnectionInterface) ([]Response, error) {

    users := map[string]uaa.User{"no-guid-yet": uaa.User{Emails: []string{options.To}}}
    space := ""
    organization := ""

    templates, err := recipe.templateLoader.Load(options.Subject, guid, clientID, options.KindID)
    if err != nil {
        return []Response{}, TemplateLoadError("An email template could not be loaded")
    }

    return recipe.courier.Mailer().Deliver(conn, templates, users, options, space, organization, clientID), nil
}

func (recipe EmailRecipe) Trim(responses []byte) []byte {
    t := Trimmer{}
    return t.TrimFields(responses, RecipientsFieldName)
}

type UAARecipe struct {
    courier CourierInterface
}

func NewUAARecipe(courier CourierInterface) UAARecipe {
    return UAARecipe{courier: courier}
}

func (recipe UAARecipe) DeliverMail(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
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
