package postal

import "github.com/cloudfoundry-incubator/notifications/models"

type Courier struct {
    tokenLoader    TokenLoader
    userLoader     UserLoader
    spaceLoader    SpaceLoader
    templateLoader TemplateLoader
    mailer         Mailer
    receiptsRepo   models.ReceiptsRepoInterface
}

type CourierInterface interface {
    Dispatch(string, TypedGUID, Options, models.ConnectionInterface) ([]Response, error)
    Mailer() MailerInterface
}

func NewCourier(tokenLoader TokenLoader, userLoader UserLoader, spaceLoader SpaceLoader,
    templateLoader TemplateLoader, mailer Mailer, receiptsRepo models.ReceiptsRepoInterface) Courier {
    return Courier{
        tokenLoader:    tokenLoader,
        userLoader:     userLoader,
        spaceLoader:    spaceLoader,
        templateLoader: templateLoader,
        mailer:         mailer,
        receiptsRepo:   receiptsRepo,
    }
}

func (courier Courier) Dispatch(clientID string, guid TypedGUID, options Options, conn models.ConnectionInterface) ([]Response, error) {
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

    templates, err := courier.templateLoader.Load(options.Subject, guid, clientID, options.KindID)
    if err != nil {
        return responses, TemplateLoadError("An email template could not be loaded")
    }

    var userGUIDs []string
    for key := range users {
        userGUIDs = append(userGUIDs, key)
    }

    err = courier.receiptsRepo.CreateReceipts(conn, userGUIDs, clientID, options.KindID)
    if err != nil {
        return responses, err
    }

    responses = courier.mailer.Deliver(conn, templates, users, options, space, organization, clientID)

    return responses, nil
}

func (courier Courier) Mailer() MailerInterface {
    return courier.mailer
}
