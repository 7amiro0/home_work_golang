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
	httpUrl = "http://calendar:8080/"

	httpAdd          = "add?userName=%v&title=%v&description=%v&notify=%v&start=%v&end=%v"
	httpList         = "list?userName=%v"
	httpListUpcoming = "listUpcoming?userName=%v&until=%v"
	httpUpdate       = "update?id=%v&title=%v&description=%v&notify=%v&start=%v&end=%v"
	httpDelete       = "delete?id=%v"
)

var (
	start, _ = time.ParseInLocation(time.RFC3339Nano, "2000-01-01T15:00:00+00:00", time.UTC)
	end, _   = time.ParseInLocation(time.RFC3339Nano, "2001-01-01T13:00:00+00:00", time.UTC)
)

func TestHTTPServer(t *testing.T) {
	//30 second because I wait to start http server
	time.Sleep(time.Second * 30)

	t.Run("http server add and list events by all time", func(t *testing.T) {
		integration_test.IncrementUserID()
		event := integration_test.CreateEvent(
			"add",
			"http",
			"justEventToAdd",
			0,
			start,
			end,
		)
		expected := storage.SliceEvents{Events: []storage.Event{event}}
		integration_test.AddEvent(t, httpUrl, httpAdd, event)

		realResult := integration_test.ListAll[storage.SliceEvents](t, httpUrl, httpList, event.User.Name)
		require.Equalf(t, expected, realResult, "expected %v but get %v", expected, realResult)
	})

	t.Run("http server update events", func(t *testing.T) {
		integration_test.IncrementUserID()
		event := integration_test.CreateEvent(
			"old",
			"httpUpdate",
			"justOldEvent",
			0,
			start,
			end,
		)

		integration_test.AddEvent(t, httpUrl, httpAdd, event)

		integration_test.IncrementUserID()
		updated := integration_test.CreateEvent(
			"UPDATED",
			"httpJunk",
			"UPDATED",
			30,
			start.Add(time.Minute*30),
			end.Add(time.Minute*30),
		)

		integration_test.AddEvent(t, httpUrl, httpAdd, updated)

		updated.ID--
		updated.User.ID--
		updated.User.Name = event.User.Name
		expected := storage.SliceEvents{Events: []storage.Event{updated}}

		integration_test.UpdateEvents(t, httpUrl, httpUpdate, updated)

		realResult := integration_test.ListAll[storage.SliceEvents](t, httpUrl, httpList, event.User.Name)
		require.Equalf(t, expected, realResult, "expected %v but get %v", expected, realResult)
	})

	t.Run("http server delete events", func(t *testing.T) {
		integration_test.IncrementUserID()
		event := integration_test.CreateEvent(
			"delete",
			"httpDelete",
			"justEventToDelete",
			0,
			start,
			end,
		)

		integration_test.AddEvent(t, httpUrl, httpAdd, event)
		var expected = storage.SliceEvents{Events: []storage.Event{}}

		resp, err := http.Get(httpUrl + fmt.Sprintf(httpDelete, event.ID))
		require.NoErrorf(t, err, "expected nil but get %q", err)
		require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)
		defer resp.Body.Close()

		realEvent := integration_test.ListAll[storage.SliceEvents](t, httpUrl, httpList, event.User.Name)
		require.Equalf(t, expected, realEvent, "expected %v but get %v", expected, realEvent)
	})

	t.Run("http server list by time", func(t *testing.T) {
		now := time.Now().UTC().Round(time.Minute)

		integration_test.IncrementUserID()
		eventOnDay := integration_test.CreateEvent(
			"day",
			"httpListByTime",
			"justEventToList",
			0,
			now.Add(10*time.Minute),
			now.Add(20*time.Minute),
		)

		eventOnWeek := integration_test.CreateEvent(
			"week",
			"httpListByTime",
			"justEventToList",
			0,
			now.Add(-20*time.Minute).Add(storage.Week),
			now.Add(-15*time.Minute).Add(storage.Week),
		)

		eventOnMonth := integration_test.CreateEvent(
			"month",
			"httpListByTime",
			"justEventToList",
			0,
			now.Add(-20*time.Minute).Add(storage.Month),
			now.Add(-15*time.Minute).Add(storage.Month),
		)

		expected := storage.SliceEvents{
			Events: []storage.Event{eventOnMonth, eventOnWeek, eventOnDay},
		}

		integration_test.AddEvent(t, httpUrl, httpAdd, eventOnDay)
		integration_test.AddEvent(t, httpUrl, httpAdd, eventOnWeek)
		integration_test.AddEvent(t, httpUrl, httpAdd, eventOnMonth)

		realResult := integration_test.ListByTime[storage.SliceEvents](t, httpUrl, httpListUpcoming, eventOnDay.User.Name, 0)
		require.Equalf(t, expected.Events[2:], realResult.Events, "expected %v but get %v", expected.Events[2:], realResult.Events)

		realResult = integration_test.ListByTime[storage.SliceEvents](t, httpUrl, httpListUpcoming, eventOnWeek.User.Name, 1)
		require.Equalf(t, expected.Events[1:], realResult.Events, "expected %v but get %v", expected.Events[1:], realResult.Events)

		realResult = integration_test.ListByTime[storage.SliceEvents](t, httpUrl, httpListUpcoming, eventOnDay.User.Name, 2)
		require.Equalf(t, expected, realResult, "expected %v but get %v", expected, realResult.Events)
	})
}
