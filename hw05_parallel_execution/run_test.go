package hw05parallelexecution

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func generateErrorTasks(tasks *[]Task, tasksCount int, runTasksCount *int32) {
	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		*tasks = append(*tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(runTasksCount, 1)
			return err
		})
	}
}

func generateNilTasks(tasks *[]Task, tasksCount int, runTasksCount *int32, learTime *time.Duration) {
	for i := 0; i < tasksCount; i++ {
		*learTime += time.Millisecond * time.Duration(rand.Intn(50))

		*tasks = append(*tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
			atomic.AddInt32(runTasksCount, 1)
			return nil
		})
	}
}

func generateNilAndErrorTasks(tasks *[]Task, tasksCount int, runTasksCount *int32, leadTime *time.Duration) {
	for i := 0; i < tasksCount/2; i++ {
		*leadTime += time.Millisecond * time.Duration(rand.Intn(100))

		*tasks = append(*tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(runTasksCount, 1)
			return nil
		})
	}
	for i := 0; i < tasksCount/2; i++ {
		*leadTime += time.Millisecond * time.Duration(rand.Intn(100))

		*tasks = append(*tasks, func() error {
			err := fmt.Errorf("error from task %d", i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(runTasksCount, 1)
			return err
		})
	}
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		generateErrorTasks(&tasks, tasksCount, &runTasksCount)
		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		leadTime := time.Duration(time.Millisecond)
		workersCount := 5
		maxErrorsCount := 1

		generateNilTasks(&tasks, tasksCount, &runTasksCount, &leadTime)

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)
			return true
		}, time.Duration(leadTime/2), time.Duration(int(leadTime)/tasksCount),
			"goruntime complete as concurrency")
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("max errors equal 0", func(t *testing.T) {
		tasksCount := 20
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		leadTime := time.Duration(time.Millisecond)
		workersCount := 10
		maxErrorsCount := 0

		generateNilAndErrorTasks(&tasks, tasksCount, &runTasksCount, &leadTime)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)
			return true
		}, time.Duration(leadTime/2), time.Duration(int(leadTime)/tasksCount),
			"goruntime complete as concurrency")
		require.Equal(t, runTasksCount/2, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("worker equal 0", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		leadTime := time.Duration(time.Millisecond)
		workersCount := 0
		maxErrorsCount := 5

		generateNilAndErrorTasks(&tasks, tasksCount, &runTasksCount, &leadTime)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)
			return true
		}, time.Duration(leadTime/2), time.Duration(int(leadTime)/tasksCount),
			"goruntime complete as concurrency")
	})
}
