package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
)

var ErrNotExistID = errors.New("%v id not exist")

type Storage struct {
	storage map[int]storage.Event
	mu      sync.RWMutex
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func New() *Storage {
	return &Storage{
		storage: make(map[int]storage.Event),
		mu:      sync.RWMutex{},
	}
}

func (s *Storage) Add(_ context.Context, event storage.Event) (err error) {
	s.mu.Lock()

	s.storage[event.ID] = event

	s.mu.Unlock()

	return err
}

func (s *Storage) Delete(_ context.Context, id int) (err error) {
	s.mu.Lock()
	if _, ok := s.storage[id]; ok {
		delete(s.storage, id)
	} else {
		err = fmt.Errorf(ErrNotExistID.Error(), id)
	}
	s.mu.Unlock()

	return err
}

func (s *Storage) Update(_ context.Context, event storage.Event) (err error) {
	s.mu.Lock()
	if _, ok := s.storage[event.ID]; ok {
		s.storage[event.ID] = event
	} else {
		err = fmt.Errorf(ErrNotExistID.Error(), event.ID)
	}
	s.mu.Unlock()

	return err
}

func (s *Storage) List(_ context.Context, idUser int) (result []storage.Event) {
	for _, event := range s.storage {
		if event.UserID == idUser {
			result = append(result, event)
		}
	}

	return result
}
