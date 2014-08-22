package postal

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

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

type DeliveryQueue struct {
    deliveries chan Delivery
}

func NewDeliveryQueue() *DeliveryQueue {
    queue := DeliveryQueue{
        deliveries: make(chan Delivery),
    }

    return &queue
}

func (queue *DeliveryQueue) Enqueue(delivery Delivery) {
    go func(delivery Delivery) {
        queue.deliveries <- delivery
    }(delivery)
}

func (queue *DeliveryQueue) Dequeue() <-chan Delivery {
    return queue.deliveries
}
