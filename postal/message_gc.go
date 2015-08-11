package postal

import (
	"log"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
)

type MessageGC struct {
	messagesRepo    messagesRepoInterface
	db              db.DatabaseInterface
	lifetime        time.Duration
	logger          *log.Logger
	timer           <-chan time.Time
	pollingInterval time.Duration
}

type messagesRepoInterface interface {
	DeleteBefore(db.ConnectionInterface, time.Time) (int, error)
}

func NewMessageGC(lifetime time.Duration, db db.DatabaseInterface,
	messagesRepo messagesRepoInterface, pollingInterval time.Duration, logger *log.Logger) MessageGC {
	return MessageGC{
		messagesRepo:    messagesRepo,
		db:              db,
		lifetime:        lifetime,
		logger:          logger,
		pollingInterval: pollingInterval,
		timer:           time.After(0),
	}
}

func (gc MessageGC) Collect() {
	threshold := time.Now().Add(-1 * gc.lifetime)
	_, err := gc.messagesRepo.DeleteBefore(gc.db.Connection(), threshold)
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
