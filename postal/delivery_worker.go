package postal

import (
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
)

const (
    StatusUnavailable = "unavailable"
    StatusFailed      = "failed"
    StatusDelivered   = "delivered"
    StatusNotFound    = "notfound"
    StatusNoAddress   = "noaddress"
)

type DeliveryWorker struct {
    guidGenerator GUIDGenerationFunc
    logger        *log.Logger
    mailClient    mail.ClientInterface
    queue         *DeliveryQueue
    halt          chan bool
}

func NewDeliveryWorker(guidGenerator GUIDGenerationFunc, logger *log.Logger, mailClient mail.ClientInterface, queue *DeliveryQueue) DeliveryWorker {
    return DeliveryWorker{
        guidGenerator: guidGenerator,
        logger:        logger,
        mailClient:    mailClient,
        queue:         queue,
        halt:          make(chan bool),
    }
}

func (worker DeliveryWorker) Run() {
    go worker.Work()
}

func (worker DeliveryWorker) Work() {
    for {
        select {
        case delivery := <-worker.queue.Dequeue():
            delivery.Response <- worker.Deliver(delivery)
        case <-worker.halt:
            return
        }
    }
}

func (worker DeliveryWorker) Halt() {
    worker.halt <- true
}

func (worker DeliveryWorker) Deliver(delivery Delivery) Response {
    var status, notificationID string

    if len(delivery.User.Emails) > 0 && strings.Contains(delivery.User.Emails[0], "@") {
        context, message := worker.Pack(delivery)
        status = worker.SendMail(message)
        notificationID = context.MessageID
    } else {
        if delivery.User.ID == "" {
            status = StatusNotFound
        } else {
            status = StatusNoAddress
        }
    }

    return Response{
        Status:         status,
        Recipient:      delivery.UserGUID,
        NotificationID: notificationID,
    }
}

func (worker DeliveryWorker) Pack(delivery Delivery) (MessageContext, mail.Message) {
    env := config.NewEnvironment()
    context := NewMessageContext(delivery.User.Emails[0], delivery.Options, env, delivery.Space, delivery.Organization, delivery.ClientID, worker.guidGenerator, delivery.Templates)
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
