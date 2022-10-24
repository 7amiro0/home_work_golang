package storage

import "time"

type Event struct {
	ID          int64
	Title       string
	UserID      int64
	Description string
	End         time.Time
	Start       time.Time
}
