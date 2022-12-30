package integration_tests

import (
	"context"
	"encoding/json"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/integration_test"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/sender"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	time.Sleep(time.Minute * 1)
	t.Run("scheduler clear old events", func(t *testing.T) {
		integration_test.IncrementUserID()
		var event storage.Event
		expected := storage.SliceEvents{Events: []storage.Event{}}

		for i := 0; i < 10; i++ {
			event = integration_test.CreateEvent(
				"junk",
				"clearEvent",
				"justEventToClear",
				0,
				start,
				end,
			)

			integration_test.AddEvent(t, grpcUrl, grpcAdd, event)
		}

		time.Sleep(3 * time.Minute)

		realEvent := integration_test.ListAll[integration_test.SliceEvents](t, grpcUrl, grpcList, event.User.Name)
		require.Equalf(t, expected, realEvent.ConvertToEvents(), "expected %v but get %v", expected, realEvent)
	})

	t.Run("main function api", func(t *testing.T) {
		loggerLevel := os.Getenv("LOGGER_LEVEL")
		logg := logger.New(loggerLevel)
		sender := s.New(logg)

		consumer, err := sender.Start(
			os.Getenv("NAME_TEST_Q"),
			os.Getenv("URL"),
			s.OptionsReadQueue{},
		)
		require.NoErrorf(t, err, "expected nil but get %q", err)

		now := time.Now().UTC().Round(time.Minute)

		integration_test.IncrementUserID()
		eventWithNotifyMinute := integration_test.CreateEvent(
			"queue",
			"queue",
			"justEventToQueue",
			1,
			now.Add(time.Minute*2),
			now.Add(time.Minute*3),
		)

		eventWithNotifyHour := integration_test.CreateEvent(
			"queue",
			"queue",
			"justEventToQueue",
			61,
			now.Add(time.Minute*62),
			now.Add(time.Minute*63),
		)

		integration_test.AddEvent(t, httpUrl, httpAdd, eventWithNotifyMinute)
		integration_test.AddEvent(t, httpUrl, httpAdd, eventWithNotifyHour)

		ctxForMinute, cancelForMinute := context.WithCancel(context.Background())
		ctxForHour, cancelForHour := context.WithCancel(context.Background())

		go func() {
			var event storage.Event
			for msg := range consumer {
				err = json.Unmarshal(msg.Body, &event)
				require.NoErrorf(t, err, "expected nil but get %q", err)

				switch event {
				case eventWithNotifyMinute:
					cancelForMinute()
				case eventWithNotifyHour:
					cancelForHour()
				default:
					require.Fail(t, "scheduler send wrong message", event)
				}

				err = msg.Ack(false)
				require.NoErrorf(t, err, "expected nil but get %q", err)
			}
		}()

		timeOutMinute := time.NewTimer(3 * time.Minute)
		defer timeOutMinute.Stop()
		select {
		case <-ctxForMinute.Done():
		case <-timeOutMinute.C:
			require.Fail(t, "time out scheduler dont send message")
		}

		timeOutHour := time.NewTimer(3 * time.Minute)
		defer timeOutHour.Stop()
		select {
		case <-ctxForHour.Done():
		case <-timeOutHour.C:
			require.Fail(t, "time out scheduler dont send message")
		}
	})
}
