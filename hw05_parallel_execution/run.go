package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func handleErrors(chQuit chan error, chErr chan error, numTask int, maxErrorNumber int) error {
	fmt.Println()
	var resError error = nil
	counterErr := 0
	counterFinish := 0
	if maxErrorNumber <= 0 {
		for _ = range chErr {

			counterFinish++

			if counterFinish == numTask {
				close(chQuit)
				close(chErr)
				break
			}
		}
	} else {
		for err := range chErr {

			counterFinish++

			if counterFinish == numTask {
				close(chQuit)
				close(chErr)
				break
			}

			if err == nil {
				continue
			}

			counterErr++

			if counterErr == maxErrorNumber {
				close(chQuit)
				resError = ErrErrorsLimitExceeded
				break
			}
		}
	}
	return resError
}

func supplyTask(chTask chan Task, tasks []Task) {
	defer close(chTask)
	for _, task := range tasks {
		chTask <- task
	}
}

func worker(chQuit <-chan error, chErr chan<- error, chTask <-chan Task, num int) {
	for {
		select {
		case <-chQuit:
			return
		default:
		}
		select {
		case <-chQuit:
			return
		case task, isOpenChTask := <-chTask:
			if isOpenChTask {
				chErr <- task()
			}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	wg := &sync.WaitGroup{}
	chTask := make(chan Task, len(tasks))
	chErr := make(chan error, n)
	chQuit := make(chan error)
	if n <= 0 {
		return nil
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			worker(chQuit, chErr, chTask, i)
		}(i)
	}

	supplyTask(chTask, tasks)
	err := handleErrors(chQuit, chErr, len(tasks), m)
	wg.Wait()
	return err
}
