package gobble

type Heartbeater struct {
	queue    QueueInterface
	ticker   TickerInterface
	haltChan chan struct{}
}

func NewHeartbeater(queue QueueInterface, ticker TickerInterface) Heartbeater {
	return Heartbeater{
		queue:    queue,
		ticker:   ticker,
		haltChan: make(chan struct{}),
	}
}

func (beater Heartbeater) Beat(job *Job) {
	beater.ticker.Start()
	for {
		select {
		case job.ActiveAt = <-beater.ticker.Tick():
			beater.queue.Requeue(job)
		case <-beater.haltChan:
			beater.ticker.Stop()
			return
		}
	}
}

func (beater Heartbeater) Halt() {
	beater.haltChan <- struct{}{}
}
