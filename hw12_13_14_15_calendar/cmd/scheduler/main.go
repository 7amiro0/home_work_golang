package main

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/scheduler"
	sqlStorage "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage/sql"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func addInQueue(schedule *s.Scheduler, exchange string, mandatory, immediate bool) error {
	var serializedEvent []byte

	events, err := schedule.List(time.Minute * 1)
	if err != nil {
		return err
	}

	for _, event := range events {
		serializedEvent, err = yaml.Marshal(event)

		err = schedule.AddInQueue(serializedEvent, exchange, mandatory, immediate)
		if err != nil {
			schedule.Logger.Error("[ERR] While add in queue: ", err)
			return err
		}

		schedule.Logger.Info("[INFO] Add event ", event, " with error ", err)
	}

	return err
}

func main() {
	config := NewConfig()
	logg := logger.New(config.Logger.Level)
	storage := sqlStorage.New(logg)

	var ticker *time.Ticker
	var mutex = &sync.Mutex{}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	schedule := s.New(ctx, storage, logg)
	logg.Info("[INFO] Created new scheduler")
	err := schedule.Start(
		config.RabbetInfo.Name,
		config.RabbetInfo.Url,
		s.OptionsQueue{
			NoWait:     config.RabbetInfo.NoWait,
			Durable:    config.RabbetInfo.Durable,
			Exclusive:  config.RabbetInfo.Exclusive,
			AutoDelete: config.RabbetInfo.AutoDelete,
		},
	)

	if err != nil {
		cancel()
		logg.Error("[ERR] While starting scheduler: ", err)
		os.Exit(1)
	}

	go func(mx *sync.Mutex) {
		logg.Info("[INFO] Start clearer")
		ticker = time.NewTicker(time.Hour * 24)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mx.Lock()
				err = schedule.Clear()
				mx.Unlock()
				if err != nil {
					logg.Error("[ERR] While working cleaner: ", err)
				}
			}
		}
	}(mutex)

	go func(mx *sync.Mutex) {
		logg.Info("[INFO] Start adder")
		ticker = time.NewTicker(time.Minute * 1)
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
		logg.Error("[ERR] while stoping scheduler: ", err)
	}

	logg.Info("[INFO] Scheduler stoped")
}