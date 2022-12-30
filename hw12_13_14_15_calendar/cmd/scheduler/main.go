package main

import (
	"context"
	"encoding/json"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/scheduler"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	addTime   time.Duration
	clearTime time.Duration
)

func addInQueue(schedule *s.Scheduler, exchange string, mandatory, immediate bool) error {
	var serializedEvent []byte

	sliceEvents, err := schedule.List(addTime)
	if err != nil {
		return err
	}

	for _, event := range sliceEvents.Events {
		serializedEvent, err = json.Marshal(event)

		err = schedule.AddInQueue(serializedEvent, exchange, mandatory, immediate)
		if err != nil {
			schedule.Logger.Error("[ERR] While add in queue: ", err)
			return err
		}

		schedule.Logger.Info("[INFO] Add event ", event)
	}

	return err
}

func main() {
	config := NewConfig()
	logg := logger.New(config.Logger.Level)
	store := storage.New(config.Logger.Level)

	addTime = config.Ticker.Add
	clearTime = config.Ticker.Clear

	var ticker *time.Ticker
	var mutex = &sync.Mutex{}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	schedule := s.New(ctx, store, logg)

	logg.Info("[INFO] Created new scheduler")

	err := schedule.Start(
		config.RabbitInfo.Name,
		config.RabbitInfo.Url,
		s.OptionsQueue{
			NoWait:     config.RabbitInfo.NoWait,
			Durable:    config.RabbitInfo.Durable,
			Exclusive:  config.RabbitInfo.Exclusive,
			AutoDelete: config.RabbitInfo.AutoDelete,
		},
	)

	if err != nil {
		cancel()
		logg.Error("[ERR] Failed to start scheduler: ", err)
		os.Exit(1)
	}

	go func(mx *sync.Mutex) {
		logg.Info("[INFO] Start cleaner")
		ticker = time.NewTicker(clearTime)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mx.Lock()
				err = schedule.Clear()
				mx.Unlock()
				if err != nil {
					logg.Error("[ERR] Cleaner: ", err)
				}
			}
		}
	}(mutex)

	go func(mx *sync.Mutex) {
		logg.Info("[INFO] Start adder")
		ticker = time.NewTicker(addTime)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mx.Lock()
				err = addInQueue(schedule, config.OptAdd.exchange, config.OptAdd.mandatory, config.OptAdd.immediate)
				mx.Unlock()
				if err != nil {
					logg.Error("[ERR] While working adder: ", err)
				}
			}
		}
	}(mutex)

	logg.Info("[INFO] Scheduler running")

	<-ctx.Done()

	if err = schedule.Stop(); err != nil {
		logg.Error("[ERR] Failed to stop scheduler: ", err)
	}

	logg.Info("[INFO] Scheduler has been stopped")
}
