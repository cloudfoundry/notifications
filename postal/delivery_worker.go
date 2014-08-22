package postal

import (
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
)

type DeliveryWorker struct {
    logger     *log.Logger
    mailClient mail.ClientInterface
    queue      *DeliveryQueue
    halt       chan bool
}

func NewDeliveryWorker(logger *log.Logger, mailClient mail.ClientInterface, queue *DeliveryQueue) DeliveryWorker {
    return DeliveryWorker{
        logger:     logger,
        mailClient: mailClient,
        queue:      queue,
        halt:       make(chan bool),
    }
}

func (worker DeliveryWorker) Run() {
    go worker.Work()
}

func (worker DeliveryWorker) Work() {
    for {
        select {
        case delivery := <-worker.queue.Dequeue():
            worker.Deliver(delivery)
        case <-worker.halt:
            return
        }
    }
}

func (worker DeliveryWorker) Halt() {
    worker.halt <- true
}

func (worker DeliveryWorker) Deliver(delivery Delivery) {
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
