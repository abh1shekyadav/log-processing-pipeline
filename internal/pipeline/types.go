package pipeline

import "time"

type Log struct {
	ID        int
	Message   string
	Processed bool
	Timestamp time.Time
}
