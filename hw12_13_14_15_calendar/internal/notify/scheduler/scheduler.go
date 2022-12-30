package scheduler

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
	"time"
)

type Scheduler struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	ctx     context.Context
	storage Storage
	Logger  Logger
}

type Storage interface {
	app.StorageScheduler
}

type Logger interface {
	app.Logger
}

type OptionsQueue struct {
	NoWait, Durable, Exclusive, AutoDelete bool
}

func New(ctx context.Context, storage Storage, logger Logger) *Scheduler {
	return &Scheduler{
		ctx:     ctx,
		storage: storage,
		Logger:  logger,
	}
}

func (s *Scheduler) AddInQueue(msg []byte, exchange string, mandatory, immediate bool) error {
	err := s.channel.Publish(
		exchange,
		s.queue.Name,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)

	return err
}

func (s *Scheduler) List(duration time.Duration) (storage.SliceEvents, error) {
	return s.storage.ListByNotify(s.ctx, duration)
}

func (s *Scheduler) Clear() error {
	return s.storage.Clear(s.ctx)
}

func (s *Scheduler) Start(name string, url string, opt OptionsQueue) (err error) {
	if err = s.storage.Connect(s.ctx); err != nil {
		s.Logger.Error("[ERR] Error connect to storage: ", err)
		return err
	}

	s.conn, err = amqp.Dial(url)
	if err != nil {
		s.Logger.Error("[ERR] Error connect to amqp: ", err)
		return err
	}

	if s.channel, err = s.conn.Channel(); err != nil {
		s.Logger.Error("[ERR] Error connect to channel: ", err)
		return err
	}

	s.queue, err = s.channel.QueueDeclare(name, opt.Durable, opt.AutoDelete, opt.Exclusive, opt.NoWait, nil)

	return err
}

func (s *Scheduler) Stop() error {
	if err := s.storage.Close(s.ctx); err != nil {
		s.Logger.Error("[ERR] Error close to storage: ", err)
		return err
	}

	if err := s.channel.Close(); err != nil {
		s.Logger.Error("[ERR] Error close to channel: ", err)
		return err
	}

	return s.conn.Close()
}
