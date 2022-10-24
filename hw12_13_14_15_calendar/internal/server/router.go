package server

import (
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	TimeFormat = "2006-01-02T15:04:05.999999999Z"

	// Query arguments.
	argID          = "id"
	argUserID      = "userID"
	argTitle       = "title"
	argDescription = "description"
	argStart       = "start"
	argEnd         = "end"
)

type EventString struct {
	id          string
	userID      string
	title       string
	description string
	start       string
	end         string
}

func newEventString(r *http.Request) EventString {
	return EventString{
		id:          r.FormValue(argID),
		userID:      r.FormValue(argUserID),
		title:       r.FormValue(argTitle),
		description: r.FormValue(argDescription),
		start:       r.FormValue(argStart),
		end:         r.FormValue(argEnd),
	}
}

func (s *Server) Add(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	fmt.Printf("%#+v", eventS)
	userID, err := strconv.ParseInt(eventS.userID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err := time.ParseInLocation(TimeFormat, eventS.start, time.UTC)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	end, err := time.ParseInLocation(TimeFormat, eventS.end, time.UTC)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	event := storage.Event{
		Title:       eventS.title,
		UserID:      userID,
		Description: eventS.description,
		End:         end,
		Start:       start,
	}

	if err = s.App.Add(s.ctx, event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	id, err := strconv.ParseInt(eventS.id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.App.Delete(s.ctx, id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	id, err := strconv.ParseInt(eventS.userID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(eventS.userID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err := time.ParseInLocation(TimeFormat, eventS.start, time.UTC)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.Logger.Error(err)
		return
	}
	end, err := time.ParseInLocation(TimeFormat, eventS.end, time.UTC)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.Logger.Error(err)
		return
	}

	newEvent := storage.Event{
		ID:          id,
		Title:       eventS.title,
		UserID:      userID,
		Description: eventS.description,
		End:         end,
		Start:       start,
	}

	if err = s.App.Update(s.ctx, newEvent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	userID, err := strconv.ParseInt(eventS.userID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := s.App.List(s.ctx, userID)

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})
	for _, event := range res {
		str := fmt.Sprintf(
			"%v %s %s %s %s\n",
			event.ID, event.Title, event.Description, event.Start, event.End)

		_, err = w.Write([]byte(str))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
