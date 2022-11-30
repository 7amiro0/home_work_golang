package memorystorage

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	event "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Add and listing event", func(t *testing.T) {
		storage := New()
		startEvent := time.Now()
		endEvent := startEvent.Add(5 * time.Minute)

		expectedResult := make([]event.Event, 0, 10)
		var testEvent event.Event

		for i := 1; i <= 10; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      int32(i),
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			expectedResult = append(expectedResult, testEvent)
			err := storage.Add(context.Background(), testEvent)
			require.Nilf(t, err, "Error: expected nil, but get %q", err)
		}

		realResult, err := storage.List(context.Background(), "user")
		require.Nilf(t, err, "Error: expected nil, but get %q", err)

		sort.Slice(realResult, func(i, j int) bool {
			return realResult[i].ID < realResult[j].ID
		})

		require.Equal(t, expectedResult, realResult)
	})

	t.Run("Delete event", func(t *testing.T) {
		storage := New()

		startEvent := time.Now()
		endEvent := startEvent.Add(5 * time.Minute)

		expectedFullResult := make([]event.Event, 0, 10)
		var testEvent event.Event

		for i := 1; i <= 10; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      int32(i),
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			expectedFullResult = append(expectedFullResult, testEvent)
			_ = storage.Add(context.Background(), testEvent)
		}

		expectedResult := append(expectedFullResult[:3], expectedFullResult[6:]...)

		err := storage.Delete(context.Background(), 4)
		require.Nilf(t, err, "Error: expected nil, but get %q", err)
		_ = storage.Delete(context.Background(), 5)
		_ = storage.Delete(context.Background(), 6)

		realResult, err := storage.List(context.Background(), "user")
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[j].ID > realResult[i].ID
		})

		require.Equal(t, expectedResult, realResult)

		err = storage.Delete(context.Background(), 4)

		require.Equal(t, fmt.Errorf(ErrNotExistID.Error(), 4), err)
	})

	t.Run("Update event", func(t *testing.T) {
		storage := New()
		startEvent := time.Now()
		endEvent := startEvent.Add(5 * time.Minute)
		expectedResult := make([]event.Event, 0, 10)

		var testEvent event.Event

		newEvent := event.Event{
			ID:    7,
			Title: "update data",
			User: event.User{
				Name: "user",
				ID:   1,
			},
			Description: "update",
			Start:       startEvent,
			End:         endEvent,
		}

		for i := 1; i <= 10; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      int32(i),
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			expectedResult = append(expectedResult, testEvent)
			_ = storage.Add(context.Background(), testEvent)
		}

		expectedResult[6] = newEvent
		sort.Slice(expectedResult, func(i, j int) bool {
			return expectedResult[j].ID > expectedResult[i].ID
		})

		err := storage.Update(context.Background(), newEvent)
		require.Nilf(t, err, "Error: expected nil, but get %q", err)

		realResult, err := storage.List(context.Background(), "user")
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[j].ID > realResult[i].ID
		})
	})

	t.Run("list by notify", func(t *testing.T) {
		storage := New()

		var testEvent event.Event
		var expectedResult []event.Event

		for i := 1; i <= 10; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      1,
				Description: "test data",
				Start:       time.Now().Add(time.Minute),
				End:         time.Now().Add(time.Minute),
			}

			expectedResult = append(expectedResult, testEvent)
			_ = storage.Add(context.Background(), testEvent)
		}

		realResult, err := storage.ListByNotify(context.Background(), 10)
		require.Nilf(t, err, "Error: expected nil, but get %q", err)
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[j].ID > realResult[i].ID
		})

		require.Equal(t, expectedResult, realResult)
	})

	t.Run("clear old events", func(t *testing.T) {
		storage := New()

		var testEvent event.Event
		expectedResult := make([]event.Event, 0, 0)

		startEvent, _ := time.Parse(time.RFC3339Nano, time.UnixDate)
		endEvent := startEvent.Add(1 * time.Minute)

		for i := 1; i <= 10; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      1,
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			_ = storage.Add(context.Background(), testEvent)
		}

		err := storage.Clear(context.Background())
		require.Nilf(t, err, "Error: expected nil, but get %q", err)

		realResult, _ := storage.List(context.Background(), "user")
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[j].ID > realResult[i].ID
		})

		require.Equal(t, expectedResult, realResult)
	})

	t.Run("data race add", func(t *testing.T) {
		storage := New()
		wg := &sync.WaitGroup{}

		startEvent := time.Now()
		endEvent := startEvent.Add(5 * time.Minute)

		expectedResult := make([]event.Event, 0, 50)
		var testEvent event.Event

		for i := 1; i <= 50; i++ {
			wg.Add(1)
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      int32(i),
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			expectedResult = append(expectedResult, testEvent)

			go func(eventTest event.Event, wg *sync.WaitGroup) {
				err := storage.Add(context.Background(), eventTest)
				require.Nilf(t, err, "Error: expected nil, but get %q", err)
				wg.Done()
			}(testEvent, wg)
		}
		wg.Wait()

		realResult, err := storage.List(context.Background(), "user")
		require.Nilf(t, err, "Error: expected nil, but get %q", err)
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[j].ID > realResult[i].ID
		})
		require.Equal(t, expectedResult, realResult)
	})

	t.Run("data race delete", func(t *testing.T) {
		storage := New()
		wg := &sync.WaitGroup{}

		startEvent := time.Now()
		endEvent := startEvent.Add(5 * time.Minute)

		expectedResult := make([]event.Event, 0, 50)
		var testEvent event.Event

		for i := 1; i <= 50; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      int32(i),
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			expectedResult = append(expectedResult, testEvent)
			_ = storage.Add(context.Background(), testEvent)
		}

		expectedResult = expectedResult[25:]

		for i := 1; i <= 25; i++ {
			wg.Add(1)
			go func(id int64, wg *sync.WaitGroup) {
				err := storage.Delete(context.Background(), id)
				require.Nilf(t, err, "Error: expected nil, but get %q", err)
				wg.Done()
			}(int64(i), wg)
		}
		wg.Wait()

		realResult, _ := storage.List(context.Background(), "user")
		sort.Slice(realResult, func(i, j int) bool {
			return realResult[j].ID > realResult[i].ID
		})
		require.Equal(t, expectedResult, realResult)
	})

	t.Run("data race list", func(t *testing.T) {
		storage := New()
		wg := &sync.WaitGroup{}

		startEvent := time.Now()
		endEvent := startEvent.Add(5 * time.Minute)

		expectedResult := make([]event.Event, 0, 50)
		var testEvent event.Event

		for i := 1; i <= 50; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      int32(i),
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			expectedResult = append(expectedResult, testEvent)
			_ = storage.Add(context.Background(), testEvent)
		}

		for i := 1; i <= 50; i++ {
			wg.Add(1)
			go func(id int64, wg *sync.WaitGroup) {
				result, err := storage.List(context.Background(), "user")
				require.Nilf(t, err, "Error: expected nil, but get %q", err)
				sort.Slice(result, func(i, j int) bool {
					return result[j].ID > result[i].ID
				})

				require.Equal(t, expectedResult, result)

				wg.Done()
			}(int64(i), wg)
		}

		wg.Wait()
	})

	t.Run("data race update", func(t *testing.T) {
		storage := New()
		wg := &sync.WaitGroup{}

		startEvent := time.Now()
		endEvent := startEvent.Add(5 * time.Minute)

		expectedResult := make([]event.Event, 0, 50)

		for i := 1; i <= 50; i++ {
			testEvent := event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      int32(i),
				Description: "test data",
				Start:       startEvent,
				End:         endEvent,
			}

			expectedResult = append(expectedResult, testEvent)
			_ = storage.Add(context.Background(), testEvent)
		}

		for i := 1; i <= 50; i++ {
			wg.Add(1)
			go func(id int64, wg *sync.WaitGroup) {
				testEvent := event.Event{
					ID:    id,
					Title: "test",
					User: event.User{
						Name: "user",
						ID:   id,
					},
					Notify:      int32(id),
					Description: "test data",
					Start:       startEvent,
					End:         endEvent,
				}

				_ = storage.Update(context.Background(), testEvent)

				wg.Done()
			}(int64(i), wg)
		}

		wg.Wait()

		result, err := storage.List(context.Background(), "user")
		require.Nilf(t, err, "Error: expected nil, but get %q", err)
		sort.Slice(result, func(i, j int) bool {
			return result[j].ID > result[i].ID
		})
		require.Equal(t, expectedResult, result)

		for i := 1; i <= 50; i++ {
			wg.Add(1)
			go func(id int64, wg *sync.WaitGroup) {
				testEvent := event.Event{
					ID:    1,
					Title: "update test",
					User: event.User{
						Name: "user",
						ID:   id,
					},
					Notify:      5,
					Description: "test data",
					Start:       startEvent,
					End:         endEvent,
				}

				err := storage.Update(context.Background(), testEvent)
				require.Nilf(t, err, "Error: expected nil, but get %q", err)
				wg.Done()
			}(int64(i), wg)
		}

		wg.Wait()
	})

	t.Run("list by notify data race", func(t *testing.T) {
		storage := New()
		wg := &sync.WaitGroup{}

		expectedResult := make([]event.Event, 0, 50)
		var testEvent event.Event

		for i := 1; i <= 50; i++ {
			testEvent = event.Event{
				ID:    int64(i),
				Title: "test",
				User: event.User{
					Name: "user",
					ID:   int64(i),
				},
				Notify:      1,
				Description: "test data",
				Start:       time.Now().Add(time.Minute),
				End:         time.Now().Add(time.Minute),
			}

			expectedResult = append(expectedResult, testEvent)
			_ = storage.Add(context.Background(), testEvent)
		}

		for i := 1; i <= 50; i++ {
			wg.Add(1)
			go func(id int64, wg *sync.WaitGroup) {
				result, err := storage.ListByNotify(context.Background(), 10)
				require.Nilf(t, err, "Error: expected nil, but get %q", err)
				sort.Slice(result, func(i, j int) bool {
					return result[j].ID > result[i].ID
				})

				require.Equal(t, expectedResult, result)

				wg.Done()
			}(int64(i), wg)
		}

		wg.Wait()
	})
}
