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
    tokenLoader    TokenLoader
    userLoader     UserLoader
    spaceLoader    SpaceLoader
    templateLoader TemplateLoader
    mailer         Mailer
}

type CourierInterface interface {
    Dispatch(string, string, NotificationType, Options) ([]Response, error)
}

func NewCourier(tokenLoader TokenLoader, userLoader UserLoader, spaceLoader SpaceLoader, templateLoader TemplateLoader, mailer Mailer) Courier {
    return Courier{
        tokenLoader:    tokenLoader,
        userLoader:     userLoader,
        spaceLoader:    spaceLoader,
        templateLoader: templateLoader,
        mailer:         mailer,
    }
}

func (courier Courier) Dispatch(rawToken, guid string, notificationType NotificationType, options Options) ([]Response, error) {
    responses := []Response{}

    token, err := courier.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    users, err := courier.userLoader.Load(notificationType, guid, token)
    if err != nil {
        return responses, err
    }

    space, organization, err := courier.spaceLoader.Load(guid, token, notificationType)
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
