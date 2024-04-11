package gobble

import "time"

type Config struct {
	WaitMaxDuration time.Duration
	MaxQueueLength  int
}
