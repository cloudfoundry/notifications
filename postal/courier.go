package postal

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/dgrijalva/jwt-go"
)

type Courier struct {
    tokenLoader    TokenLoader
    userLoader     UserLoader
    spaceLoader    SpaceLoader
    templateLoader TemplateLoader
    mailer         Mailer
}

type CourierInterface interface {
    Dispatch(string, TypedGUID, Options) ([]Response, error)
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

func (courier Courier) Dispatch(rawToken string, guid TypedGUID, options Options) ([]Response, error) {
    responses := []Response{}

    token, err := courier.tokenLoader.Load()
    if err != nil {
        return responses, err
    }

    space, organization, err := courier.spaceLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    users, err := courier.userLoader.Load(guid, token)
    if err != nil {
        return responses, err
    }

    clientToken, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    clientID := clientToken.Claims["client_id"].(string)

    templates, err := courier.templateLoader.Load(options.Subject, guid, clientID, options.Kind)
    if err != nil {
        return responses, TemplateLoadError("An email template could not be loaded")
    }

    responses = courier.mailer.Deliver(templates, users, options, space, organization, clientID)

    return responses, nil
}
