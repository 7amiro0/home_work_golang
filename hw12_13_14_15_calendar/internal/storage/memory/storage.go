package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrNotExistID    = errors.New("%v id not exist")
	ErrEventNotExist = errors.New("event with this %v id dont exist")
)

type Storage struct {
	storage sync.Map
}

func (s *Storage) ListByNotify(_ context.Context, until time.Duration) ([]storage.Event, error) {
	now := time.Now()
	end := now.Add(until)
	result := make([]storage.Event, 0, 1)

	s.storage.Range(func(key, value any) bool {
		notify := value.(storage.Event).Start.Add(-time.Duration(value.(storage.Event).Notify) * time.Minute)
		if notify.Minute() <= now.Minute() && notify.Before(end) {
			result = append(result, value.(storage.Event))
		}

		return true
	})

	return result, nil
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func New() *Storage {
	return &Storage{
		storage: sync.Map{},
	}
}

func (s *Storage) Add(_ context.Context, event *storage.Event) (err error) {
	s.storage.Store(event.ID, *event)
	return err
}

func (s *Storage) Delete(_ context.Context, id int64) (err error) {
	_, ok := s.storage.LoadAndDelete(id)
	if !ok {
		err = fmt.Errorf(ErrNotExistID.Error(), id)
	}

	return err
}

func (s *Storage) Update(_ context.Context, event *storage.Event) (err error) {
	_, ok := s.storage.Load(event.ID)
	if !ok {
		return fmt.Errorf(ErrEventNotExist.Error(), event.ID)
	}

	s.storage.Store(event.ID, *event)

	return err
}

func (s *Storage) Clear(_ context.Context) error {
	lastYear := time.Now().Add(time.Hour * 24 * 30 * 12)
	s.storage.Range(func(key, value any) bool {
		if value.(storage.Event).End.Before(lastYear) {
			s.storage.Delete(key.(int64))
		}

		return true
	})

	return nil
}

func (s *Storage) ListUpcoming(_ context.Context, userName string, until time.Duration) ([]storage.Event, error) {
	var valueEvent storage.Event
	now := time.Now().UTC()
	end := now.Add(until)
	result := make([]storage.Event, 0, 1)
	s.storage.Range(func(key, value any) bool {
		valueEvent = value.(storage.Event)
		if valueEvent.User.Name == userName {
			if valueEvent.Start.Round(storage.Day).After(now) && valueEvent.Start.Round(storage.Day).Before(end) {
				result = append(result, value.(storage.Event))
			}
		}

		return true
	})

	return result, nil
}

func (s *Storage) List(_ context.Context, userName string) ([]storage.Event, error) {
	result := make([]storage.Event, 0, 1)
	s.storage.Range(func(key, value any) bool {
		if value.(storage.Event).User.Name == userName {
			result = append(result, value.(storage.Event))
		}

		return true
	})

	return result, nil
}
