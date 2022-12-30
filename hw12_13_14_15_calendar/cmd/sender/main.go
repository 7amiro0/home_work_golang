package main

import (
	"context"
	"encoding/json"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/sender"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := NewConfig()
	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

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

			logg.Info("Hello, your event ", event.Title, " will be started in ", event.Notify, " minutes!")
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
