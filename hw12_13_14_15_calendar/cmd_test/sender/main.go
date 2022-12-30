package main

// I copy cmd sender for integration tests

import (
	"context"
	"encoding/json"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	sh "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/scheduler"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/sender"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := NewConfig()
	logg := logger.New(config.Logger.Level)
	store := storage.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	scheduler := sh.New(ctx, store, logg)
	err := scheduler.Start(
		config.RabbitInfo.NameTest,
		config.RabbitInfo.Url,
		sh.OptionsQueue{
			NoWait:     config.RabbitInfo.NoWait,
			Durable:    config.RabbitInfo.Durable,
			Exclusive:  config.RabbitInfo.Exclusive,
			AutoDelete: config.RabbitInfo.AutoDelete,
		})
	if err != nil {
		logg.Error("[ERR] Error connect to channel: ", err)
		os.Exit(1)
	}
	defer scheduler.Stop()

	sender := s.New(logg)
	logg.Info("[INFO] Created new sender")
	msgs, err := sender.Start(
		config.RabbitInfo.Name,
		config.RabbitInfo.Url,
		s.OptionsReadQueue{
			NoWait:    config.RabbitInfo.NoWait,
			AutoAck:   config.RabbitInfo.AutoAck,
			Exclusive: config.RabbitInfo.Exclusive,
			NoLocal:   config.RabbitInfo.NoLocal,
		},
	)

	if err != nil {
		cancel()
		logg.Error("[ERR] Failed to start sender: ", err)
		os.Exit(1)
	}

	go func() {
		logg.Info("[INFO] Start read queue")
		var event storage.Event
		for msg := range msgs {
			err = json.Unmarshal(msg.Body, &event)
			if err != nil {
				logg.Error("[ERR] Can`t unmarshal msg: ", err)
				continue
			}

			err = scheduler.AddInQueue(msg.Body, "", false, false)
			if err != nil {
				logg.Error("[ERR] Can`t publish msg: ", err)
			}

			err = msg.Ack(false)
			if err != nil {
				logg.Error("[ERR] Can`t ack msg: ", err)
			}
		}
	}()

	logg.Info("[INFO] Sender running")

	<-ctx.Done()

	if err = sender.Stop(); err != nil {
		logg.Error("[ERR] Failed to stop sender: ", err)
	}

	logg.Info("[INFO] Sender has been stopped")
}
