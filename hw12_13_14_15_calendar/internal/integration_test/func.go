package integration_test

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"time"
)

var (
	userID  int64 = 0
	eventID int64 = 0
)

func IncrementUserID() {
	userID++
}

func incrementEventID() {
	eventID++
}

func CreateEvent(title, user, description string, notify int32, start, end time.Time) storage.Event {
	incrementEventID()
	return storage.Event{
		End:   end,
		Start: start,
		User: storage.User{
			Name: user,
			ID:   userID,
		},
		Title:       title,
		Description: description,
		ID:          eventID,
		Notify:      notify,
	}
}
