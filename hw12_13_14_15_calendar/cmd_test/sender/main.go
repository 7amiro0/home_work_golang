package main

// I copy cmd sender for integration tests

import (
	"context"
	"encoding/json"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	s "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/notify/sender"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
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

	conn, err := amqp.Dial(config.RabbitInfo.Url)
	if err != nil {
		logg.Fatal("[FATAL] Error connect to amqp: ", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		logg.Fatal("[FATAL] Error connect to channel: ", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		config.RabbitInfo.NameTest,
		config.RabbitInfo.Durable,
		config.RabbitInfo.AutoDelete,
		config.RabbitInfo.Exclusive,
		config.RabbitInfo.NoWait,
		nil,
	)
	if err != nil {
		logg.Fatal("[FATAL] Error queue declare: ", err)
	}

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

	if err != nil {
		cancel()
		logg.Error("[ERR] While start sender: ", err)
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

			err = channel.Publish(
				config.RabbitInfo.Exchange,
				queue.Name,
				config.RabbitInfo.Mandatory,
				config.RabbitInfo.Immediate,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg.Body,
				},
			)
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
		logg.Error("[ERR] While stoping sender: ", err)
	}

	logg.Info("[INFO] Sender stoped")
}
