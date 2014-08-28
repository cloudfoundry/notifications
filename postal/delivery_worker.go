package postal

import (
    "log"
    "math"
    "strings"
    "time"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/models"
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
    logger           *log.Logger
    mailClient       mail.ClientInterface
    unsubscribesRepo models.UnsubscribesRepoInterface
    gobble.Worker
}

func NewDeliveryWorker(id int, logger *log.Logger, mailClient mail.ClientInterface, queue gobble.QueueInterface, unsubscribesRepo models.UnsubscribesRepoInterface) DeliveryWorker {
    worker := DeliveryWorker{
        logger:           logger,
        mailClient:       mailClient,
        unsubscribesRepo: unsubscribesRepo,
    }
    worker.Worker = gobble.NewWorker(id, queue, worker.Deliver)

    return worker
}

func (worker DeliveryWorker) Deliver(job *gobble.Job) {
    var delivery Delivery
    err := job.Unmarshal(&delivery)
    if err != nil {
        panic(err)
    }

    if worker.ShouldDeliver(delivery) {
        _, message := worker.Pack(delivery)
        status := worker.SendMail(message)
        if status != StatusDelivered && job.RetryCount < 10 {
            duration := time.Duration(int64(math.Pow(2, float64(job.RetryCount))))
            job.Retry(duration * time.Minute)
        }
    }
}

func (worker DeliveryWorker) ShouldDeliver(delivery Delivery) bool {
    _, err := worker.unsubscribesRepo.Find(models.Database().Connection, delivery.ClientID, delivery.Options.KindID, delivery.UserGUID)
    if err != nil {
        if (err == models.ErrRecordNotFound{}) {
            return len(delivery.User.Emails) > 0 && strings.Contains(delivery.User.Emails[0], "@")
        }
        return false
    }
    return false
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
