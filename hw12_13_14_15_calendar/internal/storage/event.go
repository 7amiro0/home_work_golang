package storage

import "time"

const (
// ListFormat = "2006-01-02"
)

type User struct {
	Name string
	ID   int64
}

type Event struct {
	End         time.Time
	Start       time.Time
	User        User
	Title       string
	Description string
	ID          int64
	Notify      int32
}

func (e Event) GetNotifyTime() time.Time {
	return e.Start.Add(-time.Minute * time.Duration(e.Notify))
}
