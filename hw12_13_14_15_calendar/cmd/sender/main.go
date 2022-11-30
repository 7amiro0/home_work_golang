package main

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/sender"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := NewConfig()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.New(config.Logger.Level)

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
		logg.Error("[ERR] While start sender: ", err)
		os.Exit(1)
	}

	go func() {
		logg.Info("[INFO] Start read queue")
		var event storage.Event
		for msg := range msgs {
			err = yaml.Unmarshal(msg.Body, &event)
			if err != nil {
				logg.Error("[ERR] Can`t unmarshal msg: ", err)
				continue
			}

			logg.Info("hello, your event will be started in ", event.Notify, " minutes!")
			err = msg.Ack(false)
			if err != nil {
				logg.Error("[ERR] Can`t ack msg: ", err)
			}
		}
	}()

	logg.Info("[INFO] Sender running")

	<-ctx.Done()

	if err = sender.Stop(); err != nil {
		logg.Error("[ERR] While stoping sender: ", err)
	}

	logg.Info("[INFO] Sender stoped")
}
