package postal

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/dgrijalva/jwt-go"
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type Response struct {
    Status         string `json:"status"`
    Recipient      string `json:"recipient"`
    NotificationID string `json:"notification_id"`
}

type NotificationType int

const (
    StatusNotFound  = "notfound"
    StatusNoAddress = "noaddress"

    IsSpace NotificationType = iota
    IsUser
)

type GUIDGenerationFunc func() (*uuid.UUID, error)

type UAAInterface interface {
    uaa.GetClientTokenInterface
    uaa.SetTokenInterface
    uaa.UsersByIDsInterface
}

type Options struct {
    ReplyTo           string
    Subject           string
    KindDescription   string
    SourceDescription string
    Text              string
    HTML              string
    Kind              string
}

type Courier struct {
    uaaClient UAAInterface

    userLoader     UserLoader
    spaceLoader    SpaceLoader
    templateLoader TemplateLoader
    mailer         Mailer
}

type CourierInterface interface {
    Dispatch(string, string, NotificationType, Options) ([]Response, error)
}

func NewCourier(uaaClient UAAInterface, userLoader UserLoader, spaceLoader SpaceLoader, templateLoader TemplateLoader, mailer Mailer) Courier {
    return Courier{
        uaaClient:      uaaClient,
        userLoader:     userLoader,
        spaceLoader:    spaceLoader,
        templateLoader: templateLoader,
        mailer:         mailer,
    }
}

func (courier Courier) Dispatch(rawToken, guid string, notificationType NotificationType, options Options) ([]Response, error) {
    responses := []Response{}

    token, err := courier.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }
    courier.uaaClient.SetToken(token.Access)

    users, err := courier.userLoader.Load(notificationType, guid, token.Access)
    if err != nil {
        return responses, err
    }

    space, organization, err := courier.spaceLoader.Load(guid, token.Access, notificationType)
    if err != nil {
        return responses, CCDownError("Cloud Controller is unavailable")
    }

    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    clientID := clientToken.Claims["client_id"].(string)

    templates, err := courier.templateLoader.Load(options.Subject, notificationType)
    if err != nil {
        return responses, TemplateLoadError("An email template could not be loaded")
    }

    responses = courier.mailer.Deliver(templates, users, options, space, organization, clientID)

    return responses, nil
}
