package memorystorage

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	event "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Add and listing event", func(t *testing.T) {
		storage := New()
		errTime := time.Now()

		expectedResult := []event.Event{}
		var testEvent event.Event
		for i := 1; i <= 10; i++ {
			if i == 5 {
				testEvent = event.Event{
					ID:          i,
					Title:       "test",
					UserID:      234,
					Description: "test data",
					End:         errTime,
					Start:       errTime,
				}
			} else {
				testEvent = event.Event{
					ID:          i,
					Title:       "test",
					UserID:      234,
					Description: "test data",
					End:         time.Now(),
					Start:       time.Now(),
				}
			}

			expectedResult = append(expectedResult, testEvent)
			err := storage.Add(context.Background(), testEvent)
			require.Nilf(t, err, "error: expected nil, but get %q", err)
		}

		realResult := storage.List(context.Background(), 234)
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[i].ID < realResult[j].ID
		})

		require.Equal(t, expectedResult, realResult)
	})

	t.Run("Delete event", func(t *testing.T) {
		storage := New()

		expectedFullResult := []event.Event{}
		for i := 1; i <= 10; i++ {
			testEvent := event.Event{
				ID:          i,
				Title:       "test",
				UserID:      234,
				Description: "test data",
				End:         time.Now(),
				Start:       time.Now(),
			}

			expectedFullResult = append(expectedFullResult, testEvent)
			err := storage.Add(context.Background(), testEvent)
			require.Nilf(t, err, "error: expected nil, but get %q", err)
		}

		expectedResult := append(expectedFullResult[:3], expectedFullResult[6:]...)

		_ = storage.Delete(context.Background(), 4)
		_ = storage.Delete(context.Background(), 5)
		_ = storage.Delete(context.Background(), 6)

		realResult := storage.List(context.Background(), 234)
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[j].ID > realResult[i].ID
		})

		require.Equal(t, expectedResult, realResult)

		err := storage.Delete(context.Background(), 4)

		require.Equal(t, fmt.Errorf(ErrNotExistID.Error(), 4), err)
	})

	t.Run("Update event", func(t *testing.T) {
		storage := New()
		newEvent := event.Event{
			ID:          7,
			Title:       "update data",
			UserID:      234,
			Description: "update",
			End:         time.Now(),
			Start:       time.Now(),
		}

		expectedResult := []event.Event{}
		for i := 1; i <= 10; i++ {
			testEvent := event.Event{
				ID:          i,
				Title:       "test",
				UserID:      234,
				Description: "test data",
				End:         time.Now(),
				Start:       time.Now(),
			}

			expectedResult = append(expectedResult, testEvent)
			err := storage.Add(context.Background(), testEvent)
			require.Nilf(t, err, "error: expected nil, but get %q", err)
		}

		realOldResult := storage.List(context.Background(), 234)
		sort.Slice(realOldResult, func(i, j int) bool {
			return realOldResult[j].ID > realOldResult[i].ID
		})
		require.Equal(t, expectedResult, realOldResult)

		expectedResult[6] = newEvent
		sort.Slice(expectedResult, func(i, j int) bool {
			return expectedResult[j].ID > expectedResult[i].ID
		})
		_ = storage.Update(context.Background(), newEvent)

		realNewResult := storage.List(context.Background(), 234)
		sort.Slice(realNewResult, func(i, j int) bool {
			return realNewResult[j].ID > realNewResult[i].ID
		})
	})

	t.Run("data race test", func(t *testing.T) {
		storage := New()
		expectedAddResult := []event.Event{}
		wg := &sync.WaitGroup{}
		wg.Add(50)

		for i := 0; i < 50; i++ {
			eventTest := event.Event{
				ID:          i,
				Title:       "event",
				UserID:      234,
				Description: "test event",
				End:         time.Now(),
				Start:       time.Now(),
			}

			expectedAddResult = append(expectedAddResult, eventTest)

			go func(eventTest event.Event, wg *sync.WaitGroup) {
				err := storage.Add(context.Background(), eventTest)
				require.Nilf(t, err, "error: expected nil, but get %q", err)
				wg.Done()
			}(eventTest, wg)
		}
		wg.Wait()

		realAddResult := storage.List(context.Background(), 234)
		sort.Slice(realAddResult, func(i, j int) bool {
			return realAddResult[j].ID > realAddResult[i].ID
		})
		require.Equal(t, expectedAddResult, realAddResult)

		wg.Add(25)
		expectedDeleteResult := expectedAddResult[25:]

		for i := 0; i < 25; i++ {
			go func(id int, wg *sync.WaitGroup) {
				err := storage.Delete(context.Background(), id)
				require.Nilf(t, err, "error: expected nil, but get %q", err)
				wg.Done()
			}(i, wg)
		}
		wg.Wait()

		realDeleteResult := storage.List(context.Background(), 234)
		sort.Slice(realDeleteResult, func(i, j int) bool {
			return realDeleteResult[j].ID > realDeleteResult[i].ID
		})
		require.Equal(t, expectedDeleteResult, realDeleteResult)
	})
}
