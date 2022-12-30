package integration_tests

import (
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/integration_test"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

const (
	grpcUrl = "http://calendar:8000/"

	grpcAdd          = "add/%v/%v/%v/%v/%q/%q"
	grpcList         = "list/%v"
	grpcListUpcoming = "listUpcoming/%v/%v"
	grpcUpdate       = "update/%v/%v/%v/%v/%q/%q"
	grpcDelete       = "delete/%v"
)

func TestGRPCServer(t *testing.T) {
	//30 second because I wait to start grpc server
	time.Sleep(time.Second * 30)

	t.Run("grpc server add and list events by all time", func(t *testing.T) {
		integration_test.IncrementUserID()
		event := integration_test.CreateEvent(
			"add",
			"grpc",
			"justEventToAdd",
			0,
			start,
			end,
		)
		expected := storage.SliceEvents{Events: []storage.Event{event}}
		integration_test.AddEvent(t, grpcUrl, grpcAdd, event)

		realResult := integration_test.ListAll[integration_test.SliceEvents](t, grpcUrl, grpcList, event.User.Name)
		require.Equalf(t, expected, realResult.ConvertToEvents(), "expected %v but get %v", expected, realResult.ConvertToEvents())
	})

	t.Run("grpc server update event", func(t *testing.T) {
		integration_test.IncrementUserID()
		event := integration_test.CreateEvent(
			"old",
			"grpcUpdate",
			"justOldEvent",
			0,
			start,
			end,
		)

		integration_test.AddEvent(t, grpcUrl, grpcAdd, event)

		integration_test.IncrementUserID()
		updated := integration_test.CreateEvent(
			"UPDATED",
			"grpcJunk",
			"UPDATED",
			30,
			start.Add(time.Minute*30),
			end.Add(time.Minute*30),
		)
		integration_test.AddEvent(t, grpcUrl, grpcAdd, updated)
		updated.ID--
		updated.User.ID--
		updated.User.Name = event.User.Name
		expected := storage.SliceEvents{Events: []storage.Event{updated}}

		integration_test.UpdateEvents(t, grpcUrl, grpcUpdate, updated)

		realResult := integration_test.ListAll[integration_test.SliceEvents](t, grpcUrl, grpcList, event.User.Name)
		require.Equalf(t, expected, realResult.ConvertToEvents(), "expected %v but get %v", expected, realResult.ConvertToEvents())
	})

	t.Run("grpc server delete events", func(t *testing.T) {
		integration_test.IncrementUserID()
		event := integration_test.CreateEvent(
			"delete",
			"grpcDelete",
			"justEventToDelete",
			0,
			start,
			end,
		)

		integration_test.AddEvent(t, grpcUrl, grpcAdd, event)
		expected := storage.SliceEvents{Events: []storage.Event{}}

		resp, err := http.Get(grpcUrl + fmt.Sprintf(grpcDelete, event.ID))
		require.NoErrorf(t, err, "expected nil but get %q", err)
		require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)
		defer resp.Body.Close()

		realEvent := integration_test.ListAll[integration_test.SliceEvents](t, grpcUrl, grpcList, event.User.Name)
		require.Equalf(t, expected, realEvent.ConvertToEvents(), "expected %v but get %v", expected, realEvent.ConvertToEvents())
	})

	t.Run("grpc server list by time", func(t *testing.T) {
		now := time.Now().UTC().Round(time.Minute)

		integration_test.IncrementUserID()
		eventOnDay := integration_test.CreateEvent(
			"day",
			"grpcListByTime",
			"justEventToList",
			0,
			now.Add(10*time.Minute),
			now.Add(20*time.Minute),
		)

		eventOnWeek := integration_test.CreateEvent(
			"week",
			"grpcListByTime",
			"justEventToList",
			0,
			now.Add(-20*time.Minute).Add(storage.Week),
			now.Add(-15*time.Minute).Add(storage.Week),
		)

		eventOnMonth := integration_test.CreateEvent(
			"month",
			"grpcListByTime",
			"justEventToList",
			0,
			now.Add(-20*time.Minute).Add(storage.Month),
			now.Add(-15*time.Minute).Add(storage.Month),
		)

		expected := storage.SliceEvents{
			Events: []storage.Event{eventOnMonth, eventOnWeek, eventOnDay},
		}

		integration_test.AddEvent(t, grpcUrl, grpcAdd, eventOnDay)
		integration_test.AddEvent(t, grpcUrl, grpcAdd, eventOnWeek)
		integration_test.AddEvent(t, grpcUrl, grpcAdd, eventOnMonth)

		realResult := integration_test.ListByTime[integration_test.SliceEvents](t, grpcUrl, grpcListUpcoming, eventOnDay.User.Name, 0)
		require.Equalf(t, expected.Events[2:], realResult.ConvertToEvents().Events, "expected %v but get %v", expected.Events[2:], realResult.ConvertToEvents().Events)

		realResult = integration_test.ListByTime[integration_test.SliceEvents](t, grpcUrl, grpcListUpcoming, eventOnWeek.User.Name, 1)
		require.Equalf(t, expected.Events[1:], realResult.ConvertToEvents().Events, "expected %v but get %v", expected.Events[1:], realResult.ConvertToEvents().Events)

		realResult = integration_test.ListByTime[integration_test.SliceEvents](t, grpcUrl, grpcListUpcoming, eventOnDay.User.Name, 2)
		require.Equalf(t, expected, realResult.ConvertToEvents(), "expected %v but get %v", expected, realResult.ConvertToEvents())
	})
}
