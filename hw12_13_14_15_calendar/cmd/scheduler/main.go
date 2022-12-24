package main

import (
	"context"
	"encoding/json"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/scheduler"
	sqlStorage "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage/sql"
	"log"
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

	schedule.Logger.Info(addTime)
	events, err := schedule.List(addTime)
	if err != nil {
		return err
	}

	for _, event := range events {
		serializedEvent, err = json.Marshal(event)

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
	addTime = config.Ticker.Add
	clearTime = config.Ticker.Clear
	log.Println("Ticker", addTime, clearTime)

	var ticker *time.Ticker
	var mutex = &sync.Mutex{}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	schedule := s.New(ctx, storage, logg)
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
		logg.Error("[ERR] While starting scheduler: ", err)
		os.Exit(1)
	}
	defer func() {
		if err = schedule.Stop(); err != nil {
			logg.Fatal("[FATAL] While stoping scheduler: ", err)
		}

		logg.Info("[INFO] Scheduler stoped")
	}()

	go func(mx *sync.Mutex) {
		logg.Info("[INFO] Start clearer")
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
					logg.Error("[ERR] While working cleaner: ", err)
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
}
