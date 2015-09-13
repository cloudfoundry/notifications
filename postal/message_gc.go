package postal

import (
	"log"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type messagesDeleter interface {
	DeleteBefore(models.ConnectionInterface, time.Time) (int, error)
}

type MessageGC struct {
	messages        messagesDeleter
	db              db.DatabaseInterface
	lifetime        time.Duration
	logger          *log.Logger
	timer           <-chan time.Time
	pollingInterval time.Duration
}

func NewMessageGC(lifetime time.Duration, db db.DatabaseInterface, messages messagesDeleter, pollingInterval time.Duration, logger *log.Logger) MessageGC {
	return MessageGC{
		messages:        messages,
		db:              db,
		lifetime:        lifetime,
		logger:          logger,
		pollingInterval: pollingInterval,
		timer:           time.After(0),
	}
}

func (gc MessageGC) Collect() {
	threshold := time.Now().Add(-1 * gc.lifetime)
	_, err := gc.messages.DeleteBefore(gc.db.Connection(), threshold)
	if err != nil {
		gc.logger.Printf("MessageGC.Collect() failed: " + err.Error())
	}
}

func (gc MessageGC) Run() {
	go func() {
		for {
			<-gc.timer
			gc.Collect()
			gc.timer = time.After(gc.pollingInterval)
		}
	}()
}
