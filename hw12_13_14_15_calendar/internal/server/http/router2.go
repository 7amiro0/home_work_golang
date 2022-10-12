package internalhttp

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
)

const (
	timeFormat = "2006-01-02 15:04:05 0000 UTC"

	// Query arguments.
	argID          = "id"
	argUserID      = "userID"
	argTitle       = "title"
	argDescription = "description"
	argStart       = "start"
	argEnd         = "end"
)

type eventString struct {
	id          string
	userID      string
	title       string
	description string
	start       string
	end         string
}

func (s *Server) getRouter() *httprouter.Router {
	s.muxServer.HandlePath("GET", "/list", s.list)
	return nil
}

func newEventString(r *http.Request) eventString {
	return eventString{
		id:          r.FormValue(argID),
		userID:      r.FormValue(argUserID),
		title:       r.FormValue(argTitle),
		description: r.FormValue(argDescription),
		start:       r.FormValue(argStart),
		end:         r.FormValue(argEnd),
	}
}

func (s *Server) add(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	userID, err := strconv.Atoi(eventS.userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err := time.ParseInLocation(timeFormat, eventS.start, time.UTC)
	if err != nil {
		l.Error(err)
		return
	}

	end, err := time.ParseInLocation(timeFormat, eventS.end, time.UTC)
	if err != nil {
		l.Error(err)
		return
	}

	event := storage.Event{
		Title:       eventS.title,
		UserID:      userID,
		Description: eventS.description,
		End:         end,
		Start:       start,
	}

	if err = s.app.Add(s.ctx, event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	id, err := strconv.Atoi(eventS.id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.app.Delete(s.ctx, id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	eventS := newEventString(r)
	id, err := strconv.Atoi(eventS.id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(eventS.userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err := time.ParseInLocation(timeFormat, eventS.start, time.UTC)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.Error(err)
		return
	}
	end, err := time.ParseInLocation(timeFormat, eventS.end, time.UTC)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.Error(err)
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

	if err = s.app.Update(s.ctx, newEvent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) list(w http.ResponseWriter, r *http.Request, p map[string]string) {
	eventS := newEventString(r)
	userID, err := strconv.Atoi(eventS.userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := s.app.List(s.ctx, userID)
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
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
