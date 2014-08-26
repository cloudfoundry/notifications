package postal

import (
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type Delivery struct {
    User         uaa.User
    Options      Options
    UserGUID     string
    Space        string
    Organization string
    ClientID     string
    Templates    Templates
    MessageID    string
}

type DeliveryWorker struct {
    logger     *log.Logger
    mailClient mail.ClientInterface
    gobble.Worker
}

func NewDeliveryWorker(id int, logger *log.Logger, mailClient mail.ClientInterface, queue gobble.QueueInterface) DeliveryWorker {
    worker := DeliveryWorker{
        logger:     logger,
        mailClient: mailClient,
    }
    worker.Worker = gobble.NewWorker(id, queue, worker.Deliver)

    return worker
}

func (worker DeliveryWorker) Deliver(job gobble.Job) {
    var delivery Delivery
    err := job.Unmarshal(&delivery)
    if err != nil {
        panic(err)
    }

    if len(delivery.User.Emails) > 0 && strings.Contains(delivery.User.Emails[0], "@") {
        _, message := worker.Pack(delivery)
        worker.SendMail(message)
    }
}

func (worker DeliveryWorker) Pack(delivery Delivery) (MessageContext, mail.Message) {
    env := config.NewEnvironment()
    context := NewMessageContext(delivery.User.Emails[0], delivery.Options, env, delivery.Space, delivery.Organization, delivery.ClientID, delivery.MessageID, delivery.Templates)
    packager := NewPackager()

    message, err := packager.Pack(context)
    if err != nil {
        panic(err)
    }

    return context, message
}

func (worker DeliveryWorker) SendMail(message mail.Message) string {
    worker.logger.Printf("Sending email to %s", message.To)

    err := worker.mailClient.Connect()
    if err != nil {
        return StatusUnavailable
    }

    err = worker.mailClient.Send(message)
    if err != nil {
        return StatusFailed
    }

    worker.logger.Print(message.Data())

    return StatusDelivered
}
