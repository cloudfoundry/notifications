package postal

import (
    "log"
    "math"
    "strings"
    "time"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/metrics"
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
    kindsRepo        models.KindsRepoInterface
    gobble.Worker
}

func NewDeliveryWorker(id int, logger *log.Logger, mailClient mail.ClientInterface, queue gobble.QueueInterface,
    unsubscribesRepo models.UnsubscribesRepoInterface, kindsRepo models.KindsRepoInterface) DeliveryWorker {
    worker := DeliveryWorker{
        logger:           logger,
        mailClient:       mailClient,
        unsubscribesRepo: unsubscribesRepo,
        kindsRepo:        kindsRepo,
    }
    worker.Worker = gobble.NewWorker(id, queue, worker.Deliver)

    return worker
}

func (worker DeliveryWorker) Deliver(job *gobble.Job) {
    var delivery Delivery
    err := job.Unmarshal(&delivery)
    if err != nil {
        metrics.NewMetric("counter", map[string]interface{}{
            "name": "notifications.worker.panic.json",
        }).Log()
        worker.Retry(job)
    }

    if worker.ShouldDeliver(delivery) {
        message := worker.pack(delivery)
        status := worker.SendMail(message)
        if status != StatusDelivered {
            worker.Retry(job)
            metrics.NewMetric("counter", map[string]interface{}{
                "name": "notifications.worker.retry",
            }).Log()
        } else {
            metrics.NewMetric("counter", map[string]interface{}{
                "name": "notifications.worker.delivered",
            }).Log()
        }
    }
}

func (worker DeliveryWorker) Retry(job *gobble.Job) {
    if job.RetryCount < 10 {
        duration := time.Duration(int64(math.Pow(2, float64(job.RetryCount))))
        job.Retry(duration * time.Minute)
    }
}

func (worker DeliveryWorker) ShouldDeliver(delivery Delivery) bool {
    conn := models.Database().Connection()
    if worker.isCritical(conn, delivery.Options.KindID, delivery.ClientID) {
        return true
    }

    _, err := worker.unsubscribesRepo.Find(conn, delivery.ClientID, delivery.Options.KindID, delivery.UserGUID)
    if err != nil {
        if (err == models.ErrRecordNotFound{}) {
            return len(delivery.User.Emails) > 0 && strings.Contains(delivery.User.Emails[0], "@")
        }
        return false
    }
    return false
}

func (worker DeliveryWorker) isCritical(conn models.ConnectionInterface, kindID, clientID string) bool {
    kind, err := worker.kindsRepo.Find(conn, kindID, clientID)
    if (err == models.ErrRecordNotFound{}) {
        return false
    }

    return kind.Critical
}

func (worker DeliveryWorker) pack(delivery Delivery) mail.Message {
    env := config.NewEnvironment()
    context := NewMessageContext(delivery.User.Emails[0], delivery.Options, env, delivery.Space, delivery.Organization, delivery.ClientID, delivery.MessageID, delivery.UserGUID, delivery.Templates)
    packager := NewPackager()

    message, err := packager.Pack(context)
    if err != nil {
        panic(err)
    }

    return message
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
