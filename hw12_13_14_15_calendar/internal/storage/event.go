package storage

import "time"

type Event struct {
	ID          int
	Title       string
	UserID      int
	Description string
	End         time.Time
	Start       time.Time
}
