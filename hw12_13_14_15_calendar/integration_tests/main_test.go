package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/integration_test"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/sender"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	//Protocol
	httpUrl = "http://calendar:8080/"
	grpcUrl = "http://calendar:8000/"

	//HTTP server work
	httpAdd          = "add?userName=%v&title=%v&description=%v&notify=%v&start=%v&end=%v"
	httpList         = "list?userName=%v"
	httpListUpcoming = "listUpcoming?userName=%v&until=%v"
	httpUpdate       = "update?id=%v&title=%v&description=%v&notify=%v&start=%v&end=%v"
	httpDelete       = "delete?id=%v"

	//GRPC server work
	grpcAdd          = "add/%v/%v/%v/%v/%q/%q"
	grpcList         = "list/%v"
	grpcListUpcoming = "listUpcoming/%v/%v"
	grpcUpdate       = "update/%v/%v/%v/%v/%q/%q"
	grpcDelete       = "delete/%v"
)

var (
	userID  int64 = 1
	eventID int64 = 1
)

func increaseID(eventID, userID int64) (int64, int64) {
	eventID++
	userID++
	return eventID, userID
}

func TestAPI(t *testing.T) {
	//60 second because I wait to start all services
	time.Sleep(time.Second * 60)
	start, errS := time.ParseInLocation(time.RFC3339Nano, "2000-01-01T15:00:00+00:00", time.UTC)
	end, errE := time.ParseInLocation(time.RFC3339Nano, "2001-01-01T13:00:00+00:00", time.UTC)
	require.NoErrorf(t, errS, "expected nil but get %q", errS)
	require.NoErrorf(t, errE, "expected nil but get %q", errE)

	t.Run("servers add and list events by all time", func(t *testing.T) {
		httpEvent := storage.Event{
			End:   end,
			Start: start,
			User: storage.User{
				Name: "http",
				ID:   userID,
			},
			Title:       "title",
			Description: "some",
			ID:          eventID,
			Notify:      1,
		}
		eventID, userID = increaseID(eventID, userID)
		expectedHTTP := storage.SliceEvents{Events: []storage.Event{httpEvent}}
		integration_test.AddEvent(t, httpUrl, httpAdd, httpEvent)

		realHTTPResult := integration_test.ListAll[storage.SliceEvents](t, httpUrl, httpList, httpEvent.User.Name)
		require.Equalf(t, expectedHTTP, realHTTPResult, "expected %v but get %v", expectedHTTP, realHTTPResult)

		grpcEvent := httpEvent
		grpcEvent.ID = eventID
		grpcEvent.User.Name = "grpc"
		grpcEvent.User.ID = userID
		eventID, userID = increaseID(eventID, userID)
		expectedGRPC := storage.SliceEvents{Events: []storage.Event{grpcEvent}}
		integration_test.AddEvent(t, grpcUrl, grpcAdd, grpcEvent)

		realGRPCResult := integration_test.ListAll[integration_test.SliceStringEvents](t, grpcUrl, grpcList, grpcEvent.User.Name)
		require.Equalf(t, expectedGRPC, realGRPCResult.ConvertToSliceEvents(), "expected %v but get %v", expectedGRPC, realGRPCResult.ConvertToSliceEvents())
	})

	t.Run("servers update events", func(t *testing.T) {
		httpEvent := storage.Event{
			End:   end,
			Start: start,
			User: storage.User{
				Name: "http",
				ID:   1,
			},
			Title:       "updateTitle",
			Description: "updateSomething",
			ID:          1,
			Notify:      1,
		}
		expectedHTTP := storage.SliceEvents{Events: []storage.Event{httpEvent}}

		grpcEvent := httpEvent
		grpcEvent.ID = 2
		grpcEvent.User.ID = 2
		grpcEvent.User.Name = "grpc"
		expectedGRPC := storage.SliceEvents{Events: []storage.Event{grpcEvent}}

		integration_test.UpdateEvents(t, httpUrl, httpUpdate, httpEvent)
		integration_test.UpdateEvents(t, grpcUrl, grpcUpdate, grpcEvent)

		realHTTPResult := integration_test.ListAll[storage.SliceEvents](t, httpUrl, httpList, httpEvent.User.Name)
		require.Equalf(t, expectedHTTP, realHTTPResult, "expected %v but get %v", expectedHTTP, realHTTPResult)

		realGRPCResult := integration_test.ListAll[integration_test.SliceStringEvents](t, grpcUrl, grpcList, grpcEvent.User.Name)
		require.Equalf(t, expectedGRPC, realGRPCResult.ConvertToSliceEvents(), "expected %v but get %v", expectedGRPC, realGRPCResult.ConvertToSliceEvents())
	})

	t.Run("servers delete events", func(t *testing.T) {
		var expected = storage.SliceEvents{Events: []storage.Event{}}

		resp, err := http.Get(httpUrl + fmt.Sprintf(httpDelete, 1))
		require.NoErrorf(t, err, "expected nil but get %q", err)
		require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)
		defer resp.Body.Close()

		resp, err = http.Get(grpcUrl + fmt.Sprintf(grpcDelete, 2))
		require.NoErrorf(t, err, "expected nil but get %q", err)
		require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)
		defer resp.Body.Close()

		realEvent := integration_test.ListAll[storage.SliceEvents](t, httpUrl, httpList, "http")
		require.Equalf(t, expected, realEvent, "expected %v but get %v", expected, realEvent)

		realEvent = integration_test.ListAll[storage.SliceEvents](t, grpcUrl, grpcList, "grpc")
		require.Equalf(t, expected, realEvent, "expected %v but get %v", expected, realEvent)
	})

	t.Run("scheduler clear old events", func(t *testing.T) {
		event := storage.Event{
			End:   end,
			Start: start,
			User: storage.User{
				Name: "clearOld",
				ID:   userID,
			},
			Title:       "oldEvent",
			Description: "justOldEvent",
			Notify:      1,
		}
		expected := storage.SliceEvents{Events: []storage.Event{}}

		for i := 0; i < 10; i++ {
			integration_test.AddEvent(t, grpcUrl, grpcAdd, event)
			eventID, _ = increaseID(eventID, 0)
		}

		time.Sleep(30 * time.Second)

		realEvent := integration_test.ListAll[integration_test.SliceStringEvents](t, grpcUrl, grpcList, event.User.Name)
		require.Equalf(t, expected, realEvent.ConvertToSliceEvents(), "expected %v but get %v", expected, realEvent)
		_, userID = increaseID(0, userID)
	})

	t.Run("servers list by time", func(t *testing.T) {
		now := time.Now().UTC().Round(time.Minute)

		eventOnDay := storage.Event{
			End:   now.Add(20 * time.Minute),
			Start: now.Add(10 * time.Minute),
			User: storage.User{
				Name: "list",
				ID:   userID,
			},
			Title:       "eventToList",
			Description: "justEventToListByTime",
			ID:          eventID,
			Notify:      1,
		}
		eventID, _ = increaseID(eventID, 0)

		eventOnWeek := eventOnDay
		eventOnWeek.Start = eventOnWeek.Start.Add(storage.Week - 10*time.Minute)
		eventOnWeek.End = eventOnWeek.End.Add(storage.Week)
		eventOnWeek.ID = eventID
		eventID, _ = increaseID(eventID, 0)

		eventOnMonth := eventOnDay
		eventOnMonth.Start = eventOnMonth.Start.Add(storage.Month - 10*time.Minute)
		eventOnMonth.End = eventOnMonth.End.Add(storage.Month)
		eventOnMonth.ID = eventID
		eventID, userID = increaseID(eventID, userID)

		expected := storage.SliceEvents{Events: []storage.Event{eventOnDay}}

		integration_test.AddEvent(t, httpUrl, httpAdd, eventOnDay)
		integration_test.AddEvent(t, httpUrl, httpAdd, eventOnWeek)
		integration_test.AddEvent(t, httpUrl, httpAdd, eventOnMonth)

		realGRPCResult := integration_test.ListByTime[integration_test.SliceStringEvents](t, grpcUrl, grpcListUpcoming, eventOnDay.User.Name, 0)
		realHTTPResult := integration_test.ListByTime[storage.SliceEvents](t, httpUrl, httpListUpcoming, eventOnDay.User.Name, 0)
		require.Equalf(t, expected, realHTTPResult, "expected %v but get %v", expected, realHTTPResult)
		require.Equalf(t, expected, realGRPCResult.ConvertToSliceEvents(), "expected %v but get %v", expected, realGRPCResult)

		expected.Events = append(expected.Events, eventOnWeek)
		realGRPCResult = integration_test.ListByTime[integration_test.SliceStringEvents](t, grpcUrl, grpcListUpcoming, eventOnWeek.User.Name, 1)
		realHTTPResult = integration_test.ListByTime[storage.SliceEvents](t, httpUrl, httpListUpcoming, eventOnWeek.User.Name, 1)
		require.Equalf(t, expected, realHTTPResult, "expected %v but get %v", expected, realHTTPResult)
		require.Equalf(t, expected, realGRPCResult.ConvertToSliceEvents(), "expected %v but get %v", expected, realGRPCResult)

		expected.Events = append(expected.Events, eventOnMonth)
		realGRPCResult = integration_test.ListByTime[integration_test.SliceStringEvents](t, grpcUrl, grpcListUpcoming, eventOnDay.User.Name, 2)
		realHTTPResult = integration_test.ListByTime[storage.SliceEvents](t, httpUrl, httpListUpcoming, eventOnDay.User.Name, 2)
		require.Equalf(t, expected, realHTTPResult, "expected %v but get %v", expected, realHTTPResult)
		require.Equalf(t, expected, realGRPCResult.ConvertToSliceEvents(), "expected %v but get %v", expected, realGRPCResult)
	})

	t.Run("main function api", func(t *testing.T) {
		loggerLevel := os.Getenv("LEVEL")
		logg := logger.New(loggerLevel)
		sender := s.New(logg)

		consumer, err := sender.Start(os.Getenv("NAME_TEST_Q"), os.Getenv("URL"), s.OptionsReadQueue{})
		require.NoErrorf(t, err, "expected nil but get %q", err)

		now := time.Now().UTC().Round(time.Minute)
		eventWithNotifyMinute := storage.Event{
			End:   now.Add(4 * time.Minute),
			Start: now.Add(3 * time.Minute),
			User: storage.User{
				Name: "queue",
				ID:   userID,
			},
			Title:       "eventToQueue",
			Description: "justEventToQueue",
			ID:          eventID,
			Notify:      1,
		}
		eventID, _ = increaseID(eventID, 0)

		eventWithNotifyHour := eventWithNotifyMinute
		eventWithNotifyHour.ID = eventID
		eventWithNotifyHour.Notify = 60
		eventWithNotifyHour.End = eventWithNotifyHour.End.Add(1 * time.Hour)
		eventWithNotifyHour.Start = eventWithNotifyHour.Start.Add(1 * time.Hour)
		eventID, userID = increaseID(eventID, userID)

		integration_test.AddEvent(t, grpcUrl, grpcAdd, eventWithNotifyMinute)
		integration_test.AddEvent(t, grpcUrl, grpcAdd, eventWithNotifyHour)

		ctxForMinute, cancelForMinute := context.WithCancel(context.Background())
		ctxForHour, cancelForHour := context.WithCancel(context.Background())

		go func() {
			var event storage.Event
			for msg := range consumer {
				err = json.Unmarshal(msg.Body, &event)
				require.NoErrorf(t, err, "expected nil but get %q", err)
				err = msg.Ack(false)
				require.NoErrorf(t, err, "expected nil but get %q", err)

				switch event {
				case eventWithNotifyMinute:
					cancelForMinute()
				case eventWithNotifyHour:
					cancelForHour()
				default:
					require.Fail(t, "Scheduler send wrong message", event)
				}
			}
		}()

		timeOut := time.NewTimer(3 * time.Minute)
		select {
		case <-ctxForMinute.Done():
		case <-timeOut.C:
			require.Fail(t, "Time out: scheduler dont send message")
		}

		timeOut.Reset(3 * time.Minute)
		select {
		case <-ctxForHour.Done():
		case <-timeOut.C:
			require.Fail(t, "Time out: scheduler dont send message")
		}
	})
}
