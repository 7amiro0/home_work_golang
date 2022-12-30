package integration_test

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"strconv"
	"time"
)

type userString struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
}

type EventString struct {
	End         time.Time  `json:"End"`
	Start       time.Time  `json:"Start"`
	User        userString `json:"User"`
	Title       string     `json:"Title"`
	Description string     `json:"Description"`
	ID          string     `json:"ID"`
	Notify      int32      `json:"Notify"`
}

func (e EventString) convertToEvent() storage.Event {
	id, _ := strconv.Atoi(e.ID)
	uID, _ := strconv.Atoi(e.User.ID)
	return storage.Event{
		End:   e.End,
		Start: e.Start,
		User: storage.User{
			Name: e.User.Name,
			ID:   int64(uID),
		},
		Title:       e.Title,
		Description: e.Description,
		ID:          int64(id),
		Notify:      e.Notify,
	}
}

type SliceEvents struct {
	Events []EventString `json:"events"`
}

func (s SliceEvents) ConvertToEvents() storage.SliceEvents {
	se := storage.SliceEvents{Events: make([]storage.Event, 0, len(s.Events))}
	for _, event := range s.Events {
		se.Events = append(se.Events, event.convertToEvent())
	}

	return se
}
