package storage

import (
	"time"
)

const (
	Day   = 24 * time.Hour
	Week  = Day * 7
	Month = Day * 30
)

type SliceEvents struct {
	Events []Event `json:"events"`
}

type User struct {
	Name string `json:"Name"`
	ID   int64  `json:"ID"`
}

type Event struct {
	End         time.Time `json:"End"`
	Start       time.Time `json:"Start"`
	User        User      `json:"User"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	ID          int64     `json:"ID"`
	Notify      int32     `json:"Notify"`
}

func (e Event) GetNotifyTime() time.Time {
	return e.Start.Add(-time.Minute * time.Duration(e.Notify))
}
