package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func AddEvent(t *testing.T, server, addCommand string, event storage.Event) {
	resp, err := http.Get(server + fmt.Sprintf(addCommand,
		event.User.Name,
		event.Title,
		event.Description,
		event.Notify,
		strings.Replace(event.Start.Format(time.RFC3339Nano), "+", "%2B", 1),
		strings.Replace(event.End.Format(time.RFC3339Nano), "+", "%2B", 1),
	))
	require.NoErrorf(t, err, "expected nil but get %q", err)
	defer resp.Body.Close()
	require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)
}

func ListAll[T storage.SliceEvents | SliceEvents](t *testing.T, server, listCommand, userName string) T {
	resp, err := http.Get(server + fmt.Sprintf(listCommand, userName))
	require.NoErrorf(t, err, "expected nil but get %q", err)
	defer resp.Body.Close()
	require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoErrorf(t, err, "expected nil but get %q", err)

	var realEvent T
	err = json.Unmarshal(body, &realEvent)
	require.NoErrorf(t, err, "expected nil but get %q", err)

	return realEvent
}

func ListByTime[T storage.SliceEvents | SliceEvents](t *testing.T, server, listCommand, userName string, per uint64) T {
	resp, err := http.Get(server + fmt.Sprintf(listCommand, userName, per))
	require.NoErrorf(t, err, "expected nil but get %q", err)
	defer resp.Body.Close()
	require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoErrorf(t, err, "expected nil but get %q", err)

	var realEvent T
	err = json.Unmarshal(body, &realEvent)
	require.NoErrorf(t, err, "expected nil but get %q", err)

	return realEvent
}

func UpdateEvents(t *testing.T, server, updateCommand string, event storage.Event) {
	resp, err := http.Get(server + fmt.Sprintf(
		updateCommand,
		event.ID,
		event.Title,
		event.Description,
		event.Notify,
		strings.Replace(event.Start.Format(time.RFC3339Nano), "+", "%2B", 1),
		strings.Replace(event.End.Format(time.RFC3339Nano), "+", "%2B", 1),
	))
	require.NoErrorf(t, err, "expected nil but get %q", err)
	defer resp.Body.Close()
	require.Equalf(t, 200, resp.StatusCode, "expected 200 but status code %v", resp.StatusCode)
}
