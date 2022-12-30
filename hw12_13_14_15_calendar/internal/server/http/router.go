package internalhttp

import (
	"encoding/json"
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"net/http"
	"strconv"
	"time"
)

const (
	// Query arguments.
	argID          = "id"
	argUserID      = "userID"
	argUserName    = "userName"
	argTitle       = "title"
	argDescription = "description"
	argNotify      = "notify"
	argStart       = "start"
	argEnd         = "end"
	argUntil       = "until"
)

type EventString struct {
	id          string
	userID      string
	userName    string
	title       string
	description string
	notify      string
	start       string
	end         string
	until       string
}

func newEventString(r *http.Request) EventString {
	return EventString{
		id:          r.FormValue(argID),
		userID:      r.FormValue(argUserID),
		userName:    r.FormValue(argUserName),
		title:       r.FormValue(argTitle),
		description: r.FormValue(argDescription),
		notify:      r.FormValue(argNotify),
		start:       r.FormValue(argStart),
		end:         r.FormValue(argEnd),
		until:       r.FormValue(argUntil),
	}
}

func (s *HTTPServer) add(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)

	notify, err := strconv.ParseInt(eventS.notify, 10, 32)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse notify: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339Nano, eventS.start)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse start: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339Nano, eventS.end)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse end: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !start.Before(end) {
		s.server.Logger.Error(fmt.Sprintf(server.ErrEndBeforeStart, end, start))
		return
	}

	event := storage.Event{
		Title:       eventS.title,
		User:        storage.User{Name: eventS.userName},
		Description: eventS.description,
		Notify:      int32(notify),
		End:         end.UTC(),
		Start:       start.UTC(),
	}

	if err = s.server.App.Add(s.server.Ctx, &event); err != nil {
		s.server.Logger.Error("[ERR] Cannot add event: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *HTTPServer) delete(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	id, err := strconv.ParseInt(eventS.id, 10, 64)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.server.App.Delete(s.server.Ctx, id); err != nil {
		s.server.Logger.Error("[ERR] Cannot delete event: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *HTTPServer) update(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	id, err := strconv.ParseInt(eventS.id, 10, 64)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notify, err := strconv.ParseInt(eventS.notify, 10, 32)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse notify: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339Nano, eventS.start)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse start: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339Nano, eventS.end)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse end: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !start.Before(end) {
		s.server.Logger.Error(fmt.Sprintf(server.ErrEndBeforeStart, end, start))
		return
	}

	event := storage.Event{
		ID:          id,
		Title:       eventS.title,
		Description: eventS.description,
		Notify:      int32(notify),
		End:         end.UTC(),
		Start:       start.UTC(),
	}

	if err = s.server.App.Update(s.server.Ctx, &event); err != nil {
		s.server.Logger.Error("[ERR] Cannot update event: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *HTTPServer) listUpcoming(w http.ResponseWriter, r *http.Request) {
	var until time.Duration
	eventS := newEventString(r)

	untilUint, err := strconv.ParseUint(eventS.until, 10, 64)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot parse until: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch untilUint {
	case 0:
		until = storage.Day
	case 1:
		until = storage.Week
	case 2:
		until = storage.Month
	}

	res, err := s.server.App.ListUpcoming(s.server.Ctx, eventS.userName, until)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot list events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	str, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(str)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *HTTPServer) list(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	res, err := s.server.App.List(s.server.Ctx, eventS.userName)
	if err != nil {
		s.server.Logger.Error("[ERR] Cannot list events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	str, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(str)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
